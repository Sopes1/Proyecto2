/**
 * Required External Modules
 */
const express = require("express");
const cors = require("cors");
const redisController = require("./controllers/redis");
const mongoController = require("./controllers/mongo");
/**
 * App Variables 3
 */
const port = 8080;
const app = express();
const server = require("http").createServer(app);
const io = require("socket.io")(server, {
  cors: { methods: ["GET", "PATCH", "POST", "PUT"], origin: "*" },
});
/**
 *  App Configuration
 */
app.use(express.json());
app.use(express.urlencoded({ extended: true }));
app.use(cors());

/**
 * Routes Definitions
 */

app.get("/", (req, res) => {
  res.status(200).send("WHATABYTE: Food For Devs");
});

/**
 * Socket io
 */
let interval;
let interval2;
io.on("connection", (socket) => {
  console.log("We have a new conecction!!");
  if (interval) {
    clearInterval(interval);
    clearInterval(interval2);
  }

  socket.on("top", async () => {
    interval = setInterval(async () => {
      let ultimos = await redisController.getlast();
      socket.emit("last", ultimos);
    }, 2000);

    interval2 = setInterval(async () => {
      let top = await redisController.getall();
      socket.emit("top", top);
    }, 2000);
  });

  socket.on("status", async (id_player) => {
    interval = setInterval(async () => {
      let report = await redisController.getstatus(id_player);
      console.log(report);
      socket.emit("player", report);
    }, 2000);
  });

  socket.on("graficas", async () => {
    interval = setInterval(async () => {
      let worker = await mongoController.getWorkers();
      console.log(worker);
      socket.emit("worker", worker);
    }, 2000);
    interval2 = setInterval(async () => {
      let topgame = await mongoController.getTopGame();
      console.log(topgame);
      socket.emit("topgame", topgame);
    }, 2000);
  });

  socket.on("datamongo", async () => {
    interval = setInterval(async () => {
      let datam = await mongoController.getData();
      console.log(datam);
      socket.emit("data", datam);
    }, 2000);
    interval2 = setInterval(async () => {
      let logsm = await mongoController.getLogs();
      console.log(logsm);
      socket.emit("logs", logsm);
    }, 2000);
  });

  socket.on("insert", async () => {
    await redisController.insertdata();
  });
  socket.on("disconnect", () => {
    console.log("Client had left!!");
    clearInterval(interval);
    clearInterval(interval2);
  });
});

/**
 * Server Activation
 */
server.listen(port, () => {
  console.log(`Listening to requests on http://localhost:${port}`);
});
