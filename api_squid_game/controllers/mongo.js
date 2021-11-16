const { MongoClient } = require("mongodb");
let dotenv = require("dotenv").config();

const client = new MongoClient(process.env.mongoURL_);

client.connect();

async function getTop() {
  const result = await client
    .db("Proyecto2Sopes")
    .collection("Games")
    .aggregate([
      {
        $group: {
          _id: "$winner",
          count: { $sum: 1 },
        },
      },
      {
        $sort: { count: -1 },
      },
    ])
    .limit(2)
    .toArray();
  return result;
}
async function getLast() {
  const result = await client
    .db("Proyecto2Sopes")
    .collection("Games")
    .find({})
    .project({ _id: 0, game: 1, gamename: 1, winner: 1, players: 1 })
    .limit(2)
    .sort({ _id: -1 })
    .toArray();
  return result;
}

async function getWorkers() {
  let data = {
    label: "",
    y: 0,
  };
  let datapoint = [];
  const result = await client
    .db("Proyecto2Sopes")
    .collection("Logs")
    .aggregate([
      {
        $group: {
          _id: "$worker",
          count: { $sum: 1 },
        },
      },
    ])
    .limit(3)
    .toArray();
  console.log(result);
  result.map((row) => {
    data.label = row._id;
    data.y = row.count;
    datapoint.push(data);
  });
  return datapoint;
}
async function getTopGame() {
  let data = {
    label: "",
    y: 0,
  };
  let datapoint = [];
  const result = await client
    .db("Proyecto2Sopes")
    .collection("Logs")
    .aggregate([
      {
        $group: {
          _id: "$gamename",
          count: { $sum: 1 },
        },
      },
    ])
    .limit(3)
    .toArray();
  console.log(result);
  result.map((row) => {
    data.label = row._id;
    data.y = row.count;
    datapoint.push(data);
  });
  return datapoint;
}
async function getLogs() {
  const result = await client
    .db("Proyecto2Sopes")
    .collection("Logs")
    .find()
    .project({
      _id: 0,
      request_game: 1,
      game: 1,
      gamename: 1,
      winner: 1,
      players: 1,
      worker: 1,
    })
    .sort({ _id: -1 })
    .toArray();
  return result;
}

async function getData() {
  const result = await client
    .db("Proyecto2Sopes")
    .collection("Games")
    .find()
    .project({ _id: 0, game: 1, gamename: 1, winner: 1, players: 1 })
    .sort({ _id: -1 })
    .toArray();
  return result;
}

module.exports = {
  getTop,
  getLast,
  getWorkers,
  getLogs,
  getData,
  getTopGame,
};
