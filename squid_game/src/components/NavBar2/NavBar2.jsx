import React, { Component } from "react";
import { Link } from "react-router-dom";
import $ from "jquery";
import "./NavBar2.css";
class NavBar2 extends Component {
  timer = null;
  constructor(props) {
    super(props);
    this.state = {};
  }
  test() {
    var tabsNewAnim = $("#navbarSupportedContent");
    var activeItemNewAnim = tabsNewAnim.find(".active");
    var activeWidthNewAnimHeight = activeItemNewAnim.innerHeight();
    var activeWidthNewAnimWidth = activeItemNewAnim.innerWidth();
    var itemPosNewAnimTop = activeItemNewAnim.position();
    var itemPosNewAnimLeft = activeItemNewAnim.position();
    $(".hori-selector").css({
      top: itemPosNewAnimTop.top + "px",
      left: itemPosNewAnimLeft.left + "px",
      height: activeWidthNewAnimHeight + "px",
      width: activeWidthNewAnimWidth + "px",
    });
    $("#navbarSupportedContent").on("click", "li", function (e) {
      $("#navbarSupportedContent ul li").removeClass("active");
      $(this).addClass("active");
      var activeWidthNewAnimHeight = $(this).innerHeight();
      var activeWidthNewAnimWidth = $(this).innerWidth();
      var itemPosNewAnimTop = $(this).position();
      var itemPosNewAnimLeft = $(this).position();
      $(".hori-selector").css({
        top: itemPosNewAnimTop.top + "px",
        left: itemPosNewAnimLeft.left + "px",
        height: activeWidthNewAnimHeight + "px",
        width: activeWidthNewAnimWidth + "px",
      });
    });
  }

  componentDidMount() {
    this.test();
    window.addEventListener("resize", this.test);
  }
  componentWillUnmount() {
    window.addEventListener("resize", this.test);
  }
  logout = () => {
    console.log("salir");
  };

  render() {
    return (
      <>
        <div className="sticky-top">
          <nav className="navbar navbar-expand-custom navbar-mainbg">
            <Link className="navbar-brand navbar-logo" to="#">
              SQUID GAME
            </Link>
            <button
              className="navbar-toggler"
              type="button"
              aria-controls="navbarSupportedContent"
              aria-expanded="false"
              aria-label="Toggle navigation"
            >
              <i className="fas fa-bars"></i>
            </button>
            <div
              className="collapse navbar-collapse"
              id="navbarSupportedContent"
            >
              <ul className="navbar-nav ml-auto">
                <div className="hori-selector">
                  <div className="left"></div>
                  <div className="right"></div>
                </div>
                <li className="nav-item active">
                  <Link className="nav-link" to="/topreport">
                    <i className="fas fa-list-ol"></i>Top 10
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" to="/status">
                    <i className="fas fa-users">
                      <i className="fas fa-tachometer-alt"></i>
                    </i>
                    Status
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" to="/report">
                    <i className="far fa-chart-bar"></i>Charts
                  </Link>
                </li>
                <li className="nav-item">
                  <Link className="nav-link" to="/datamongo">
                    <i className="fas fa-list-ol"></i>Data Mongo
                  </Link>
                </li>
                <li className="nav-item">
                  <Link
                    className="nav-link"
                    to="#"
                    onClick={this.logout.bind(this)}
                  >
                    <i className="fas fa-sign-out-alt"></i>Log Out
                  </Link>
                </li>
              </ul>
            </div>
          </nav>
        </div>
      </>
    );
  }
}

export default NavBar2;
