var vm = new Vue({
  el: '#app',
  data: {
      message: 'Hello, Vue!',
      blocks: [],
      bpage: 1,
      blimit: 10,
      block: {},
      transactions: [],
      txpage: 1,
      txlimit: 10,
      selectedParticipant: ''
  },
  methods: {
      getBlocks: function (page) {
          this.bpage = page;
          axios
              .get("http://localhost:3000/blocks?select=number,hash,txcount&order=number.desc&limit=" + this.blimit
                  + "&offset=" + this.blimit * (page - 1))
              .then(response => (this.blocks = response.data));
      },
      prevBPage: function () {
          if (this.bpage > 1) {
              this.getBlocks(this.bpage - 1);
          }
      },
      nextBPage: function () {
          this.getBlocks(this.bpage + 1);
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
      getTxsByParticipant: function (hash, page) {
          this.txpage = page;
          this.selectedParticipant = hash;
          axios
              .get("http://localhost:3000/transactions?select=*&or=(from.eq." + hash + ",to.eq." + hash
                  + ")&order=blocknumber.desc&limit=" + this.txlimit + "&offset=" + this.txlimit * (page - 1))
              .then(response => (this.transactions = response.data));
      },
      prevTxPage: function () {
          this.getTxsByParticipant(this.selectedParticipant, this.txpage - 1);
      },
      nextTxPage: function () {
          this.getTxsByParticipant(this.selectedParticipant, this.txpage + 1);
      }

  }
})