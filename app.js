var vm = new Vue({
  el: '#app',
  data: {
      message: 'Hello, Vue!',
      blocks: [],
      block: {}
  },
  methods: {
      getBlocks: function () {
          axios
              .get("http://localhost:3000/blocks?select=number,hash,txcount&order=number.desc&limit=10")
              .then(response => (this.blocks = response.data));
      },
      getBlock: function (num) {
          axios
              .get("http://localhost:3000/blocks?select=*&number=eq." + num)
              .then(response => (this.block = response.data[0]));
      }
  }
})