var vm = new Vue({
  el: '#app',
  data: {
      message: 'Hello, Vue!',
      blocks: [],
      block: {},
      transactions: []
  },
  methods: {
      getBlocks: function () {
          axios
              .get("http://localhost:3000/blocks?select=number,hash,txcount&order=number.desc&limit=10&txcount=gt.5")
              .then(response => (this.blocks = response.data));
      },
      getBlock: function (num) {
          axios
              .get("http://localhost:3000/blocks?select=*&number=eq." + num)
              .then(response => (this.block = response.data[0]));
      },
      getTxByHash: function (hash) {
          axios
              .get("http://localhost:3000/transactions?select=*&hash=eq." + hash)
              .then(response => (this.transactions = response.data));
      },
      getTxsByBlockNum: function (num) {
          axios
              .get("http://localhost:3000/transactions?select=*&blocknumber=eq." + num)
              .then(response => (this.transactions = response.data));
      },
      getTxsByParticipant: function (hash) {
          axios
              .get("http://localhost:3000/transactions?select=*&or=(from.eq." + hash + ",to.eq." + hash
                  + ")&order=blocknumber.desc&limit=200")
              .then(response => (this.transactions = response.data));
      }

  }
})