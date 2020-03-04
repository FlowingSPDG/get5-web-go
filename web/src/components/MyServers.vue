<template>
  <ul class="list-group" v-if="servers">
    <li class="list-group-item" v-if="servers.length == 0">
      {{$t('MyServers.NoServersFound')}}
    </li>

    <el-table :data="servers" style="width: 100%" v-else>
      <el-table-column prop="id" label="Server ID" width="180"></el-table-column>
      <el-table-column prop="display_name" label="Display Name" width="180"></el-table-column>
      <el-table-column prop="ip_string" label="IP Address"></el-table-column>
      <el-table-column prop="port" label="Port"></el-table-column>
      <el-table-column label="Status">
        <template slot-scope="scope">
          <template v-if="scope.row.in_use"> {{$t('MyServers.InUse')}} </template>
          <template v-else> {{$t('MyServers.Free')}}  </template>
        </template>
      </el-table-column>
      <el-table-column>
        <template slot-scope="scope">
          <el-button icon="el-icon-edit" @click="$router.push('/server/'+scope.row.id+'/edit')"></el-button>
          <el-button icon="el-icon-delete" @click="DeleteTeam(scope.row.id)" v-if="!scope.row.in_use"></el-button>
        </template>
      </el-table-column>
    </el-table>
  </ul>
</template>

<script>
export default {
  name: 'MyServers',
  data () {
    return {
      servers: [],
      user: {}
    }
  },
  async created () {
    this.user = await this.CheckLoggedIn()
    let servers = await this.GetUserData(res.data.userid)
    this.servers = servers.servers
  },
  methods: {
    async DeleteTeam (serverid) {
      try {
        await this.$confirm('This will permanently delete the server. Continue?', 'Warning', {
          confirmButtonText: 'OK',
          cancelButtonText: 'Cancel',
          type: 'warning'
        })
        let res = await this.axios.delete(`/api/v1/server/${serverid}/delete`)
        this.$message({
          message: this.$t('MyServers.MessageDeleteSuccess'),
          type: 'success'
        })
        this.$router.push('/myservers')
      } catch (err) {
        if (err.response) {
          if (typeof err.response.data === 'string') {
            this.$message.error(err.response.data)
          } else if (typeof err.response.data === 'object') {
            this.$message.error(err.response.data.errormessage)
          }
        }
      }
    }
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
