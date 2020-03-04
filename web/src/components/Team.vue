<template>

  <div class="container">
    <h1 v-cloak>
      <img :src="get_flag_link(team)" /> {{ team.name }} {{ team.logo }}
      <div class="pull-right" v-if="Deletable">
        <el-button icon="el-icon-delete" @click="DeleteTeam(team.id)"></el-button>
      </div>
      <div class="pull-right" v-if="Editable">
        <el-button icon="el-icon-edit" @click="$router.push('/team/'+team.id+'/edit')"></el-button>
      </div>
    </h1>

    <br>

    <div class="panel panel-default">
      <div class="panel-heading">{{$t('Team.Players')}}</div>
      <div class="panel-body" v-cloak>
        <el-table :data="players">
          <el-table-column label="SteamID" prop="steamid" width="400">
            <template slot-scope="scope">
              <a :href="'https://steamcommunity.com/profiles/'+scope.row.steamid">{{scope.row.steamid}}</a>
            </template>
          </el-table-column>
          <el-table-column
            prop="name"
            label="Name"
            width="180">
          </el-table-column>
        </el-table>
      </div>
    </div>

    <div class="panel panel-default">
      <div class="panel-heading">{{$t('Team.RecentMatches')}}</div>
        <div class="panel-body"  v-if="team">
          <div v-for="(match, index) in matches" :key="index" >
            <router-link :to="'/match/'+match.id">#{{match.id}}</router-link>: {{ matchdata[index] }}
            <br>
          </div>
      </div>
    </div>

  </div>
</template>

<script>
export default {
  name: 'Team',
  data () {
    return {
      team: {
        flag: '',
        name: '',
        logo: ''
      },
      matches: [],
      matchdata: [],
      players: [],
      teamdatas: {},
      user: {
        isLoggedIn: false,
        steamid: 0,
        userid: 0
      },
      Editable: false,
      Deletable: false
    }
  },
  async created () {
    const teamdataPromise = this.GetTeamData(this.$route.params.teamid)
    this.matches = await this.GetRecentMatches(this.$route.params.teamid)
    this.team = await teamdataPromise
    if (!Array.isArray(this.matches)) {
      this.matches = []
    }
    for (let i = 0; i < this.matches.length; i++) {
      let res = await this.get_vs_match_result(this.matches[i])
      this.matchdata.push(res)
    }
    if (!Array.isArray(this.team.steamids)) {
      this.team.steamids = []
    }
    for (let i = 0; i < this.team.steamids.length; i++) {
      this.GetSteamName(this.team.steamids[i])
    }
    this.user = await CheckLoggedIn()
    this.Editable = this.CheckTeamEditable(this.user.userid)
    this.Deletable = this.CheckTeamDeletable(this.user.userid)
  },
  methods: {
    async GetSteamName (steamid) {
      let self = this
      if (steamid === '') {
        return
      }
      return new Promise(async (resolve, reject) => {
        const res = await this.axios.get(`/api/v1/GetSteamName?steamID=${steamid}`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.$store.state.auth}`
          },
          data: null
        })
        self.players.push({ steamid: steamid, name: res.data })
        resolve(res.data)
      })
    },
    CheckTeamEditable: function (userid) {
      return this.team.user_id === parseInt(userid)
    },
    CheckTeamDeletable: function (userid) {
      return this.team.user_id === parseInt(userid)
    },
    async get_vs_match_result (match) {
      return new Promise(async (resolve, reject) => {
        let MyScore
        let OtherTeamScore
        let OtherTeam
        const teamid = parseInt(this.$route.params.teamid)
        const maxmaps = parseInt(match.max_maps)

        if (match.team1.id === teamid) {
          MyScore = parseInt(match.team1_score)
          OtherTeamScore = parseInt(match.team2_score)
          OtherTeam = await this.GetTeamData(match.team2.id)
          // for a bo1 replace series score with the map score
          if (maxmaps === 1) {
            if (!match.map_stats) {
              match.map_stats = []
            }
            if (match.map_stats.length === 1) {
              if (match.team1_id === this.team.id) {
                MyScore = parseInt(match.map_stats[0].team1_score)
                OtherTeamScore = parseInt(match.map_stats[0].team2_score)
              } else {
                MyScore = parseInt(match.map_stats[0].team1_score)
                OtherTeamScore = parseInt(match.map_stats[0].team2_score)
              }
            }
          }
          if (match.live) {
            let r = `Live, ${MyScore}:${OtherTeamScore} vs ${OtherTeam.name}`
            resolve(r)
          }
          if (MyScore < OtherTeamScore) {
            let r = `Lost ${MyScore}:${OtherTeamScore} vs ${OtherTeam.name}`
            resolve(r)
          } else if (MyScore > OtherTeamScore) {
            let r = `Won ${MyScore}:${OtherTeamScore} vs ${OtherTeam.name}`
            resolve(r)
          } else {
            let r = `Tied ${OtherTeamScore}:${MyScore} vs ${OtherTeam.name}`
            resolve(r)
          }
        } else {
          MyScore = parseInt(match.team2_score)
          OtherTeamScore = parseInt(match.team1_score)
          OtherTeam = await this.GetTeamData(match.team1.id)
          // for a bo1 replace series score with the map score
          if (maxmaps === 1) {
            if (!Array.isArray(match.map_stats)) {
              match.map_stats = []
            }
            if (match.map_stats.length === 1) {
              if (match.team1_id === this.team.id) {
                MyScore = parseInt(match.map_stats[0].team2_score)
                OtherTeamScore = parseInt(match.map_stats[0].team1_score)
              } else {
                MyScore = parseInt(match.map_stats[0].team2_score)
                OtherTeamScore = parseInt(match.map_stats[0].team1_score)
              }
            }
          }
          if (match.live) {
            let r = `Live, ${MyScore}:${OtherTeamScore} vs ${OtherTeam.name}`
            resolve(r)
          }
          if (MyScore < OtherTeamScore) {
            let r = `Lost ${MyScore}:${OtherTeamScore} vs ${OtherTeam.name}`
            resolve(r)
          } else if (MyScore > OtherTeamScore) {
            let r = `Won ${MyScore}:${OtherTeamScore} vs ${OtherTeam.name}`
            resolve(r)
          } else {
            let r = `Tied ${OtherTeamScore}:${MyScore} vs ${OtherTeam.name}`
            resolve(r)
          }
        }
      })
    },
    async DeleteTeam (teamid) {
      try {
        await this.$confirm('This will permanently delete the team. Continue?', 'Warning', {
          confirmButtonText: 'OK',
          cancelButtonText: 'Cancel',
          type: 'warning'
        })
        let res = await this.axios.delete(`/api/v1/team/${teamid}/delete`)
        this.$message({
          message: this.$t('Team.MessageDeleteSuccess'),
          type: 'success'
        })
        this.$router.push('/myteams')
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
