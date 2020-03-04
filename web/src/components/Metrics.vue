<template>
  <div id="metrics">
    <h1>{{$t('Metrics.Get5Metrics')}}</h1>
    <el-table :data="metrics">
      <el-table-column
        prop="users"
        label="Registered Users"
        width="300">
      </el-table-column>
      <el-table-column
        prop="saved_teams"
        label="Saved Teams"
        width="300">
      </el-table-column>
      <el-table-column
        prop="matches_created"
        label="Matches created"
        width="300">
      </el-table-column>
      <el-table-column
        prop="completed_matches"
        label="Completed matches"
        width="300">
      </el-table-column>
      <el-table-column
        prop="servers_added"
        label="Servers added"
        width="300">
      </el-table-column>
      <el-table-column
        prop="maps_with_stats"
        label="Maps with stats saved"
        width="300">
      </el-table-column>
      <el-table-column
        prop="unique_players"
        label="Unique players"
        width="300">
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
export default {
  name: 'Metrics',
  data () {
    return {
      metrics: []
    }
  },
  async created () {
    const res = await this.axios.get('/api/v1/GetMetrics', {
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.$store.state.auth}`
      },
      data: null
    })
    this.metrics.push(res.data)
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
