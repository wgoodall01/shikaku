"use strict";

const tableEl = document.querySelector("#entry_table>tbody");
const rowsInput = document.querySelector(`input[name="rows"]`);
const colsInput = document.querySelector(`input[name="cols"]`);
const shrinkable = document.querySelector("#entry_shrinkable");

const max = 40;
const shrinkAt = 15;

const mkInputEl = () => {
  let td = document.createElement("td");
  let el = document.createElement("input");
  el.type = "text";
  el.pattern = "[0-9]+";
  el.dataset.hjWhitelist = "true";
  el.name = "e";
  td.appendChild(el);
  return td;
};

function update() {
  let numRows = parseInt(rowsInput.value);
  let numCols = parseInt(colsInput.value);

  if (numRows > 40) {
    numRows = max;
    rowsInput.value = max;
  }

  if (numCols > max) {
    numCols = max;
    colsInput.value = max;
  }

  if (numCols > shrinkAt || numRows > shrinkAt) {
    shrinkable.classList.add("entry_shrink");
  } else {
    shrinkable.classList.remove("entry_shrink");
  }

  console.log(`update ${numRows} x ${numCols}`);

  let rows = tableEl.querySelectorAll("tr");
  if (rows.length < numRows) {
    // Add rows to fill
    for (let i = 0; i < numRows - rows.length; i++) {
      tableEl.appendChild(document.createElement("tr"));
    }
  } else if (rows.length >= numRows) {
    // Remove rows to trim
    for (let i = numRows; i < rows.length; i++) {
      rows[i].remove();
    }
  }

  rows = tableEl.querySelectorAll("tr"); // refresh from DOM
  for (let row of rows) {
    const cols = row.querySelectorAll("td");
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

tableEl.querySelector("tr").remove();

document.addEventListener("DOMContentLoaded", () => update());
//update();
