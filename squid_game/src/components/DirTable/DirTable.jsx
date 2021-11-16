import React, { Component } from "react";
import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableContainer from "@material-ui/core/TableContainer";
import TableHead from "@material-ui/core/TableHead";
import TableRow from "@material-ui/core/TableRow";
import Typography from "@material-ui/core/Typography";
import "./DirTable.css";

class DirTable extends Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  getTitleTable() {
    return this.props.dirTitle.map((row, index) => (
      <TableCell align="center" key={index}>
        <Typography style={{ fontWeight: "bold", color: "blue" }}>
          {row}
        </Typography>
      </TableCell>
    ));
  }
  getDataTable() {
    return this.props.dirData.map((row, index) => (
      <TableRow className="rowTable" key={index}>
        {Object.values(row).map((cel, index) => (
          <TableCell align="center" key={index}>
            <Typography> {cel}</Typography>
          </TableCell>
        ))}
      </TableRow>
    ));
  }
  render() {
    return (
      <TableContainer>
        <Table aria-label="simple table">
          <TableHead>
            <TableRow>{this.getTitleTable()}</TableRow>
          </TableHead>
          <TableBody>{this.getDataTable()}</TableBody>
        </Table>
      </TableContainer>
    );
  }
}

export default DirTable;
