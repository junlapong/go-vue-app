<template>
  <div class="test">
    <h1>This is a test page</h1>
    <h2>Call REST API</h2>
    <button type="button" v-on:click="ping()">Ping</button> |
    <button type="button" v-on:click="now()">Now</button>
    <div>
      <span id="message">{{ msg }}</span>
    </div>
  </div>
</template>

<script>
import axios from "axios";

const API = axios.create({
  baseURL: `${process.env.BASE_URL}api/`,
  timeout: 1000
});

export default {
  data() {
    return {
      msg: ""
    };
  },
  methods: {
    async ping() {
      // this.msg = "Pong";
      await API.get("/ping")
        .then(res => (this.msg = res.data.message))
        .catch(err => (this.msg = err.message));
    },
    async now() {
      // this.msg = "Now";
      const data = { timestamp: new Date().getTime() };
      await API.post("/now", data)
        .then(res => (this.msg = res.data.timestamp))
        .catch(err => (this.msg = err.message));
    }
  }
};
</script>
