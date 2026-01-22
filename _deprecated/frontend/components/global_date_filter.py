import datetime
import typing

import streamlit as st

""


def get_default_range() -> typing.Tuple[datetime.datetime, datetime.datetime]:
    if "date_range" not in st.session_state:
        today = datetime.datetime.now()
        end_exclusive = today + datetime.timedelta(days=1)
        start = today.replace(day=1)
        st.session_state["date_range"] = (start, end_exclusive)
    else:
        range = st.session_state["date_range"]
        start = range[0]
        if len(range) == 2:
            end_exclusive = range[1]
        else:
            # because pandas dates have this limit. Like why??
            end_exclusive = datetime.datetime(2262, 4, 11)

    return start, end_exclusive
