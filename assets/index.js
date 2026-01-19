function filterTable(inputId, tableId) {
  // Declare variables
  let input = document.getElementById(inputId);
  let filter = input.value;
  let table = document.getElementById(tableId);
  let tr = table.getElementsByTagName("tr");

  // Loop through all table rows, and hide those who don't match the search query
  // start at 1 to skip header row
  for (i = 1; i < tr.length; i++) {
    let matched = false;
    datas = tr[i].getElementsByTagName("td");
    for (td of datas) {
      if (td) {
        txtValue = td.textContent || td.innerText;
        if (txtValue.indexOf(filter) > -1) {
          matched = true;
          break;
        }
      }
    }

    if (matched) {
      tr[i].style.display = "";
    } else {
      tr[i].style.display = "none";
    }
  }
}
