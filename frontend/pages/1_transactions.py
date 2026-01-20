import datetime
from functools import reduce

import pandas as pd
import requests
import streamlit as st

st.set_page_config(
    page_title="Transactions",
)


@st.cache_data()
def api_transactions():
    resp = requests.get("http://0.0.0.0:8080/api/transactions")
    resp.raise_for_status()

    data = resp.json().get("data", [])
    df = pd.DataFrame(data)
    print(df.dtypes)
    # df["DestinationAccount"] = (
    #     df["DestinationAccount"].astype(str).apply(lambda x: json.loads(x)["Name"])
    # )
    # df["SourceAccount"] = (
    #     df["SourceAccount"].astype(str).apply(lambda x: json.loads(x)["Name"])
    # )
    df.drop(
        inplace=True,
        columns=[
            "ID",
            "UpdatedAt",
            "CreatedAt",
            "DeletedAt",
            "SourceAccountID",
            "DestinationAccountID",
        ],
    )
    df["Date"] = pd.to_datetime(df["Date"])
    df["Amount"] = df["Amount"].astype(float)
    return df


def get_data():
    (start, _end) = st.session_state.date_filter
    end_exclusive = _end + datetime.timedelta(days=1)

    df = api_transactions()

    df = df[df["Date"] >= pd.to_datetime(start)]
    df = df[df["Date"] < pd.to_datetime(end_exclusive)]
    return df


def present(df):
    df["Amount"] = df["Amount"].apply(lambda x: format(x, ","))
    return df


def filter_data():
    text = st.session_state.transaction_filter

    if not text:
        st.session_state["data"] = df
    else:
        masks = []
        for col in df.columns:
            masks.append(df[col].astype(str).str.contains(text))

        mask = reduce(lambda x, y: x | y, masks)
        filtered = df[mask]
        st.session_state["data"] = filtered


# Date picker
today = datetime.datetime.now()
month_start = today.replace(day=1)
if "date_filter" not in st.session_state:
    st.session_state["date_filter"] = (month_start, today)

if "data" not in st.session_state:
    st.session_state["data"] = get_data()

if "transaction_filter" not in st.session_state:
    st.session_state["transaction_filter"] = ""


# TODO: implement quick filters as pills https://docs.streamlit.io/develop/api-reference/widgets/st.pills
#   also think about adding categories to these pills


st.date_input(
    "Date Range",
    (month_start, today),
    format="YYYY-MM-DD",
    key="date_filter",
    on_change=filter_data,
)

# category chart
st.write("# Amount By Category")
cats = get_data()[["Amount", "Category"]].groupby("Category").sum()
cats.drop(["Income"], inplace=True)
st.bar_chart(cats, sort=True)


st.write("# Transactions")
st.text_input(
    label="Text filter table here", on_change=filter_data, key="transaction_filter"
)

df = present(st.session_state.data)
styled = df.style.map(
    lambda x: f"color: {'red' if x.startswith('-') else 'green'}",
    subset="Amount",
)
st.dataframe(styled)
