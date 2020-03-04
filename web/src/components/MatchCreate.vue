<template>
<div id="container">
    <div id="content">
        <el-form ref="form" :model="form" :rules="rules" label-position="left" label-width="120px">
            <el-form-item :label="$t('MatchCreate.FormServer')" style="width: 653px;" prop="display_name">
                <el-select v-model="form.server_id" filterable>
                    <el-option v-for="(server, index) in servers" :label="server.display_name" :key="index" :value="server.id"></el-option>
                </el-select>
                <el-button icon="el-icon-plus" @click="$router.push('/server/create')"></el-button>

            </el-form-item>

            <el-form-item :label="$t('MatchCreate.FormTeam1')" style="width: 653px;" prop="team1_id">
                <el-select v-model="form.team1_id" filterable>
                    <el-option v-for="(team, index) in teams" :label="team.name" :key="index" :value="team.id"></el-option>
                </el-select>
                <el-button icon="el-icon-plus" @click="$router.push('/team/create')"></el-button>
            </el-form-item>

            <el-form-item v-if="match_text_option" :label="$t('MatchCreate.FormTeam1Title')" style="width: 653px;" prop="team1_string">
                <el-input v-model="form.team1_string"></el-input>
            </el-form-item>

            <el-form-item :label="$t('MatchCreate.FormTeam2')" style="width: 653px;" prop="team2_id">
                <el-select v-model="form.team2_id" filterable>
                    <el-option v-for="(team, index) in teams" :label="team.name" :key="index" :value="team.id"></el-option>
                </el-select>
                <el-button icon="el-icon-plus" @click="$router.push('/team/create')"></el-button>
            </el-form-item>

            <el-form-item v-if="match_text_option" :label="$t('MatchCreate.FormTeam2Title')" style="width: 653px;" prop="team2_string">
                <el-input v-model="form.team2_string"></el-input>
            </el-form-item>

            <el-form-item v-if="match_text_option" :label="$t('MatchCreate.FormMatchTitle')" style="width: 653px;" prop="match_title">
                <el-input v-model="form.match_title"></el-input>
            </el-form-item>

            <el-form-item :label="$t('MatchCreate.FormSeriesType')" style="width: 653px;" prop="series_type">
                <el-radio-group v-model="form.series_type" v-on:change="form.veto_mappool = []">
                    <el-radio v-for="(option, index) in series_type" :label="option.type" :key="index"></el-radio>
                </el-radio-group>
            </el-form-item>

            <el-form-item :label="$t('MatchCreate.FormMapPool')" style="width: 653px;" prop="veto_mappool">
                <el-checkbox-group v-if="form.series_type !== 'bo1-preset'" v-model="form.veto_mappool">
                    <el-checkbox v-for="(map, index) in mappool.active" :label="map" :key="index"></el-checkbox >
                </el-checkbox-group>
                <el-radio v-else-if="form.series_type === 'bo1-preset'" v-for="(map, index) in mappool.active" v-model="form.veto_mappool[0]" :label="map" :key="index"></el-radio>
            </el-form-item>

            <el-form-item style="width: 653px;">
              <el-input v-model="cvar"></el-input>
                <el-button icon="el-icon-plus" type="primary" @click="AddCvar()">{{$t('MatchCreate.AddCvar')}}</el-button>
            </el-form-item>

            <el-form-item :label="$t('MatchCreate.FormCVARS')" style="width: 653px;" prop="cvars">
              <div v-for="(cvar, index) in cvars" :key="index">
              <el-input v-model="cvars[index].value">
                <template slot="prepend">{{cvar.cvar}}</template>
                <el-button icon="el-icon-delete" slot="append" @click="DeleteCvar(index)"></el-button>
              </el-input>
              </div>
            </el-form-item>

            <el-form-item style="width: 653px;">
                <el-button type="primary" @click="RegisterMatch">{{$t('MatchCreate.CreateMatch')}}</el-button>
            </el-form-item>
        </el-form>
    </div>
</div>
</template>

<script>
export default {
  name: 'MatchCreate',
  data () {
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
      servers: [],
      teams: [],
      mappool: {
        active: ['de_dust2', 'de_mirage', 'de_inferno', 'de_overpass', 'de_train', 'de_nuke', 'de_vertigo'],
        reserve: ['de_cache', 'de_season']
      },

      /* mappool: [
        {
          system: 'de_dust2',
          formal: 'Dust II'
        },
        {
          system: 'de_mirage',
          formal: 'Mirage'
        },
        {
          system: 'de_inferno',
          formal: 'Inferno'
        },
        {
          system: 'de_overpass',
          formal: 'Overpass'
        },
        {
          system: 'de_nuke',
          formal: 'NUKE'
        },
        {
          system: 'de_train',
          formal: 'Train'
        },
        {
          system: 'de_vertigo',
          formal: 'Vertigo'
        },
        {
          system: 'de_cache',
          formal: 'Cache'
        }
      ], */
      series_type: [
        {
          type: 'bo1-preset',
          desc: 'Bo1 with preset map'
        },
        {
          type: 'bo1',
          desc: 'Bo1 with map vetoes'
        },
        {
          type: 'bo2',
          desc: 'Bo2 with map vetoes'
        },
        {
          type: 'bo3',
          desc: 'Bo3 with map vetoes'
        },
        {
          type: 'bo5',
          desc: 'Bo5 with map vetoes'
        },
        {
          type: 'bo7',
          desc: 'Bo7 with map vetoes'
        }
      ],
      cvar: '',
      cvars: [],
      form: {
        server_id: 0,
        team1_id: undefined,
        team2_id: undefined,
        max_maps: 0,
        title: '',
        skip_veto: false,
        veto_mappool: [],
        series_type: 'bo1',
        cvars: {}
      },
      rules: {
        server_id: [{
          required: true,
          trigger: 'change',
          message: 'Please chose server'
        }],
        team1_id: [{
          required: true,
          trigger: 'change',
          message: 'Please chose team1 id'
        }],
        team2_id: [{
          required: true,
          trigger: 'change',
          message: 'Please chose team2 id'
        }],
        title: [{
          required: true,
          trigger: 'change',
          message: 'Please input title'
        }],
        skip_veto: [{
          required: false,
          trigger: 'change',
          message: 'Please chose skip veto option'
        }],
        veto_mappool: [{
          required: false,
          trigger: 'change',
          message: 'Please chose map(s)'
        }]
      },
      match_text_option: false // TODO
    }
  },
  async created () {
    this.user = await this.CheckLoggedIn()
    this.servers = await this.GetServers()
    this.teams = await this.GetTeams()
    this.mappool = await this.GetMapList()
  },
  methods: {
    async RegisterMatch () {
      this.form.team1_id = parseInt(this.form.team1_id, 10)
      this.form.team2_id = parseInt(this.form.team2_id, 10)
      this.form.server_id = parseInt(this.form.server_id, 10)
      if (this.form.team1_id === this.form.team2_id) {
        this.$message.error('Teams cannot be equal')
        return
      }
      switch (this.form.series_type) {
        case 'bo1-preset':
          this.form.max_maps = 1
          this.form.skip_veto = true
          if (this.form.veto_mappool.length !== 1) {
            this.$message.error('You must have exactly 1 map selected to do a bo1 with a preset map')
            return
          }
          break
        case 'bo1':
          this.form.max_maps = 1
          break
        case 'bo2':
          this.form.max_maps = 3
          break
        case 'bo3':
          this.form.max_maps = 3
          break
        case 'bo5':
          this.form.max_maps = 5
          break
        case 'bo7':
          this.form.max_maps = 7
          break
      }
      for (let i = 0; i < this.cvars.length; i++) {
        this.form.cvars[this.cvars[i].cvar] = this.cvars[i].value
      }
      const json = JSON.stringify(this.form)
      this.$refs['form'].validate(async (valid) => {
        if (valid) {
          try {
            let res = await this.axios.post('/api/v1/match/create', json)
            this.form = {}
            this.$message({
              message: this.$t('MatchCreate.MessageRegisterSuccess'),
              type: 'success'
            })
            this.$router.push('/mymatches')
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
    async AddCvar () {
      this.cvars.push({ cvar: this.cvar, value: '' })
      this.cvar = ''
    },
    async DeleteCvar (index) {
      this.cvars.splice(index, 1)
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
    display: block;
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
