<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <p>
      This website will help you calculate any instagram account engagement rate.
    </p>
    <h3>Enter Instagram account name</h3>
    <div class="form">
      <form v-on:submit="submitForm">
        <input type="text" v-model="username" placeholder="Username">
        <button type="submit">Submit</button>
      </form>
      <div v-if="show">
        <p>Using (likes + comments) / followers</p>
        <p>Engagement Rate for
          <span class="engagement">{{ username }} </span> is
          <span class="engagement">{{ engagement }} %</span></p>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: String
  },

  methods: {
    submitForm(e) {
      e.preventDefault()
      let user = this.username
      if (user) {
        this.show = true
        this.axios.get('http://localhost:8080/api/username/' + user).then((response) => {
          this.engagement = response.data
        })
      } else {
        this.show = false
      }
    }
  },

  data: function () {
    return {
      username: "",
      engagement: 0,
      show: false
    }
  },
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.form {
  padding: 5px;
}

input[type=text], button[type=submit] {
  padding: 10px;
}

button[type=submit] {
  background-color: darkblue;
  color: white;
  border: none;
  margin-left: 10px;
}

.engagement {
  font-weight: bold;
}
</style>
