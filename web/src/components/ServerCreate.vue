<template>
<div id="container">
<div id="content">
  <el-form ref="form" :model="form" :rules="rules" label-position="left" label-width="120px">
  <el-form-item :label="$t('ServerCreate.FormDisplayName')" style="width: 653px;" prop="display_name">
    <el-input v-model="form.display_name"></el-input>
  </el-form-item>

  <el-form-item :label="$t('ServerCreate.FormServerIP')" style="width: 653px;" prop="ip_string">
    <el-input v-model="form.ip_string"></el-input>
  </el-form-item>

  <el-form-item :label="$t('ServerCreate.FormServerPort')" style="width: 653px;" prop="port">
    <el-input v-model="form.port" placeholder=27015></el-input>
  </el-form-item>

  <el-form-item :label="$t('ServerCreate.FormRCONPassword')" style="width: 653px;" prop="rcon_password">
    <el-input type="password"  v-model="form.rcon_password"></el-input>
    <p class="help-block">{{$t('ServerCreate.InfoWillNotBeExplosed')}}</p>
  </el-form-item>

  <el-form-item :label="$t('ServerCreate.FormPublicTeam')" style="width: 653px;" prop="public_server" v-if="user.admin">
    <el-switch v-model="form.public_server"></el-switch>
  </el-form-item>

  <el-form-item style="width: 653px;" v-if="edit">
    <el-button type="primary" v-if="edit" @click="UpdateServer">{{$t('misc.Update')}}</el-button>
  </el-form-item>

  <el-form-item style="width: 653px;" v-else>
    <el-button type="primary" v-if="!edit" @click="RegisterServer">{{$t('misc.Create')}}</el-button>
  </el-form-item>
</el-form>

</div>
</div>
</template>

<script>
export default {
  name: 'ServerCreate',
  props: ['edit'],
  data () {
    var ValidateIPaddress = (rule, value, callback) => {
      if (/^(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/.test(value)) {
        callback()
      } else {
        callback(new Error('Please input valid IP'))
      }
    }
    return {
      user: {
        id: 0,
        steam_id: '',
        name: '',
        admin: false,
        servers: null,
        teams: null,
        matches: null
      },
      form: {
        ip_string: '',
        port: 27015,
        rcon_password: '',
        display_name: '',
        public_server: false
      },
      rules: {
        ip_string: [
          { required: true, trigger: 'change', validator: ValidateIPaddress, message: 'Please enter valid IP' }
        ],
        port: [
          { required: true, trigger: 'change', message: 'Please enter valid Port' }
        ],
        rcon_password: [
          { required: true, message: 'Please input rcon password', trigger: 'change' }
        ],
        display_name: [
          { required: false, message: 'Please input server display name', trigger: 'change' }
        ],
        public_server: [
          { required: false, message: 'Please check if server is public', trigger: 'change' }
        ]
      }
    }
  },
  async created () {
    this.$message({
      dangerouslyUseHTMLString: true,
      message: this.$t('ServerCreate.Get5Help'),
      type: 'info',
      duration: 1000 * 10,
      showClose: true
    })
    this.user = await CheckLoggedIn()
    if (this.edit) {
      try {
        let res = await this.axios.get(`api/v1/server/${this.$route.params.serverID}/GetServerInfo`)
        this.form.ip_string = res.data.ip_string
        this.form.port = res.data.port
        this.form.rcon_password = res.data.rcon_password
        this.form.display_name = res.data.display_name
        this.form.public_server = res.data.public_server
      } catch (err) {
        if (typeof err.response.data === 'string') {
          this.$message.error(err.response.data)
        } else if (typeof err.response.data === 'object') {
          this.$message.error(err.response.data.errormessage)
        }
      }
    }
  },
  methods: {
    async RegisterServer () {
      const json = JSON.stringify(this.form)
      this.$refs['form'].validate(async (valid) => {
        if (valid) {
          try {
            let res = await this.axios.post('/api/v1/server/create', json)
            this.form = {}
            this.$message({
              message: this.$t('ServerCreate.MessageRegisterSuccess'),
              type: 'success'
            })
            this.$router.push('/matches')
          } catch (err) {
            if (typeof err.response.data === 'string') {
              this.$message.error(err.response.data)
            } else if (typeof err.response.data === 'object') {
              this.$message.error(err.response.data.errormessage)
            }
          }
        } else {
          this.$message.error('Please fill form')
        }
      })
    },
    async UpdateServer () {
      const json = JSON.stringify(this.form)
      this.$refs['form'].validate(async (valid) => {
        if (valid) {
          try {
            let res = await this.axios.put(`/api/v1/server/${this.$route.params.serverID}/edit`, json)
            this.form = {}
            this.$message({
              message: this.$t('ServerCreate.MessageeEditSuccess'),
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
        } else {
          this.$message.error('Please fill form')
        }
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->

<style scoped>
h1,
h2 {
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

#container {
  padding-left: 500px;
  margin-right: auto;
  margin-left: auto;
  align-content: center;
}
</style>
