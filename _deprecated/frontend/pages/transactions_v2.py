import streamlit as st
from components import data

st.set_page_config(page_title="Transactions", layout="wide")

DATA_DIR = "./data/csvs"


file = st.selectbox("Datasets", data.list_files(), key=data.FILES_KEY)

if file:
    trans = data.transactions(file)
    if len(trans) > 0:
        trans = trans.style.map(
            lambda x: f"color: {'red' if x < 0 else 'green'}",
            subset="Amount",
        )
    st.dataframe(trans)
