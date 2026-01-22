import glob
import os
import typing

import pandas as pd
import streamlit as st

from frontend.components import constants

st.set_page_config(page_title="Transactions", layout="wide")

TRANSACTION_DATA_KEY = "data"
FILES_KEY = "files"
CURRENT_FILE_KEY = "current_file"


def list_files() -> typing.List[str]:
    st.session_state[FILES_KEY] = sorted(
        f
        for f in glob.glob(os.path.join(constants.DATA_DIR, "**/*.csv"), recursive=True)
    )
    return st.session_state.files


def transactions(file: str) -> pd.DataFrame:
    if CURRENT_FILE_KEY not in st.session_state:
        st.session_state[CURRENT_FILE_KEY] = file
        st.session_state[TRANSACTION_DATA_KEY] = pd.read_csv(file)
    else:
        if st.session_state[CURRENT_FILE_KEY] != file:
            df = pd.read_csv(file)
            df["Amount"] = df["Amount"].roun(2)
            df.sort_values("Date", ascending=False, inplace=True)
            st.session_state[TRANSACTION_DATA_KEY] = df

    return st.session_state.data
