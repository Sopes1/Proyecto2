const redis = require("redis");
let dotenv = require("dotenv").config();

let id = 0;
const redisCliente = redis
  .createClient({
    host: process.env.redisURL_,
    port: process.env.redisPort_,
    db: 0,
  })
  .on("error", function (err) {
    console.error(
      process.env.redisURL_ + ":" + process.env.redisPort_ + " " + err
    );
  })
  .on("connect", function () {
    console.log(
      "Redis connected " + process.env.redisURL_ + ":" + process.env.redisPort_
    );
  });


async function getlast() {
  return new Promise((resolve, reject) => {
    redisCliente.lrange("Game", 0, 9, (error, reply) => {
      let datas = [];
      reply.map((row) => {
        datas.unshift(JSON.parse(row));
      });
      return error ? reject(error) : resolve(datas);
    });
  });
}
async function getall() {
  return new Promise((resolve, reject) => {
    redisCliente.lrange("Game", 0, -1, function (err, reply) {
      let ganadores = new Map();
      reply.map((row) => {
        let dato = JSON.parse(row);
        if (ganadores.has(dato.winner)) {
          let vic = ganadores.get(dato.winner);
          ganadores.delete(dato.winner);
          ganadores.set(dato.winner, vic + 1);
        } else {
          ganadores.set(dato.winner, 1);
        }
      });
      ganadores[Symbol.iterator] = function* () {
        yield* [...this.entries()].sort((a, b) => b[1] - a[1]);
      };
      let ganar = Array.from(ganadores);
      return err ? reject(err) : resolve(ganar.slice(0, 9));
    });
  });
}

async function getstatus(id_player) {
  return new Promise((resolve, reject) => {
    redisCliente.lrange("Game", 0, -1, function (err, reply) {
      let data = [];
      reply.map((row) => {
        let dato = JSON.parse(row);
        if (dato.winner == id_player) {
          data.unshift(dato);
        }
      });
      return err ? reject(err) : resolve(data);
    });
  });
}

module.exports = {
  insertdata,
  getall,
  getlast,
  getstatus,
};
