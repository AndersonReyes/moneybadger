import datetime
from functools import reduce

import pandas as pd
import requests
import streamlit as st
from components import global_date_filter

st.set_page_config(
    page_title="Transactions",
)


# @st.cache_data()
def api_transactions(
    startDate: datetime.datetime, endDate: datetime.datetime
) -> pd.DataFrame:
    resp = requests.get(
        "http://0.0.0.0:8080/api/transactions",
        params={
            "StartDate": startDate.strftime("%Y-%m-%dT%H:%M:%SZ"),
            "EndDate": endDate.strftime("%Y-%m-%dT%H:%M:%SZ"),
        },
    )
    resp.raise_for_status()

    data = resp.json().get("data", [])
    df = pd.DataFrame(data)
    if data:
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


def get_data() -> pd.DataFrame:
    start, end_exclusive = global_date_filter.get_default_range()
    df = api_transactions(start, end_exclusive)
    if len(df) > 0:
        df = df[df["Date"] >= pd.to_datetime(start, utc=True)]
        df = df[df["Date"] < pd.to_datetime(end_exclusive, utc=True)]
    return df


def present(df):
    df["Amount"] = df["Amount"].apply(lambda x: format(x, ","))
    df["Date"] = df["Date"].dt.date
    return df


def refresh_data():
    text = st.session_state.transaction_filter

    df = get_data()

    if not text:
        st.session_state["data"] = df
    else:
        masks = []
        for col in df.columns:
            masks.append(df[col].astype(str).str.contains(text))

        mask = reduce(lambda x, y: x | y, masks)
        filtered = df[mask]
        st.session_state["data"] = filtered


if "data" not in st.session_state:
    st.session_state["data"] = get_data()

if "transaction_filter" not in st.session_state:
    st.session_state["transaction_filter"] = ""


# TODO: implement quick filters as pills https://docs.streamlit.io/develop/api-reference/widgets/st.pills
#   also think about adding categories to these pills

start, end = global_date_filter.get_default_range()
st.date_input(
    "Date Range",
    (start, end),
    format="YYYY-MM-DD",
    key="date_range",
    on_change=refresh_data,
)

# category chart
st.write("# Amount By Category")
cats = get_data()
if len(cats) > 0:
    cats = cats[["Amount", "Category"]].groupby("Category").sum()
    cats.drop(["Income"], inplace=True)
st.bar_chart(cats, sort=True)


st.write("# Transactions")
st.text_input(
    label="Text filter table here", on_change=refresh_data, key="transaction_filter"
)

trans = st.session_state.data
if len(trans) > 0:
    trans = present(trans)
    trans = trans.style.map(
        lambda x: f"color: {'red' if x.startswith('-') else 'green'}",
        subset="Amount",
    )
st.dataframe(trans)
