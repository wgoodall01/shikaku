'use strict';

const tableEl = document.querySelector('#entry_table>tbody');
const rowsInput = document.querySelector(`input[name="rows"]`);
const colsInput = document.querySelector(`input[name="cols"]`);

const mkInputEl = () => {
  let td = document.createElement('td');
  let el = document.createElement('input');
  el.type = 'text';
  el.pattern = '[0-9]+';
  el.name = 'e';
  td.appendChild(el);
  return td;
};

function update() {
  const numRows = parseInt(rowsInput.value);
  const numCols = parseInt(colsInput.value);

  console.log(`update ${numRows} x ${numCols}`);

  let rows = tableEl.querySelectorAll('tr');
  if (rows.length < numRows) {
    // Add rows to fill
    for (let i = 0; i < numRows - rows.length; i++) {
      tableEl.appendChild(document.createElement('tr'));
    }
  } else if (rows.length >= numRows) {
    // Remove rows to trim
    for (let i = numRows; i < rows.length; i++) {
      rows[i].remove();
    }
  }

  rows = tableEl.querySelectorAll('tr'); // refresh from DOM
  for (let row of rows) {
    const cols = row.querySelectorAll('td');
    if (cols.length < numCols) {
      for (let i = 0; i < numCols - cols.length; i++) {
        row.appendChild(mkInputEl());
      }
    } else if (cols.length >= numCols) {
      for (let i = numCols; i < cols.length; i++) {
        cols[i].remove();
      }
    }
  }
}

rowsInput.oninput = () => update();
colsInput.oninput = () => update();

tableEl.querySelector('tr').remove();
update();
