<template>
  <div id="content" class="matches">
    <h1 v-if="my_matches">{{$t('Matches.YourMatches')}}</h1>
    <h1 v-else-if="all_matches">{{$t('Matches.AllMatches')}}</h1>
    <h1 v-else>Matches for <router-link :to="'/user/'+match_owner.id">{{match_owner.name}}</router-link></h1>

  <table class="table table-striped">
    <thead>
      <tr>
        <th>{{$t('Matches.MatchID')}}</th>
        <th>{{$t('Matches.Team1')}}</th>
        <th>{{$t('Matches.Team2')}}</th>
        <th>{{$t('Matches.Status')}}</th>
        <th v-if="my_matches">{{$t('Matches.Server')}}</th>
        <th v-if="my_matches"></th>
        <th v-else>{{$t('Matches.Owner')}}</th>
      </tr>
    </thead>
    <tbody>

      <tr v-for="(match, index) in matches" :key="index" align="left">
        <td v-if="match"><router-link :to="'/match/'+match.id">{{match.id}}</router-link></td>

        <td v-if="matchinfo[match.id]">
          <img v-if="matchinfo[match.id].team1" :src="get_flag_link(matchinfo[match.id].team1)" />
          <router-link v-if="matchinfo[match.id].team1" :to="'/team/'+match.team1_id">{{matchinfo[match.id].team1.name}}</router-link>
        </td>

        <td v-if="matchinfo[match.id]">
          <img v-if="matchinfo[match.id].team1" :src="get_flag_link(matchinfo[match.id].team2)" />
          <router-link v-if="matchinfo[match.id].team2" :to="'/team/'+match.team2_id">{{matchinfo[match.id].team2.name}}</router-link>
        </td>

        <td v-if="matchinfo[match.id]">
          {{ matchinfo[match.id].status }}
        </td>

        <td v-if="my_matches && matchinfo[match.id]"><div v-if="matchinfo[match.id].server">{{ matchinfo[match.id].server.display_name }}</div> </td>
        <td v-if="my_matches && matchinfo[match.id]">
          <a v-if="(match.pending || match.live)" :href="'/match/'+match.id+'cancel'" class="btn btn-danger btn-xs align-right">Cancel</a>
        </td>
        <td v-if="!my_matches && matchinfo[match.id]">
          <router-link :to="'/user/='+matchinfo[match.id].user.id">{{ matchinfo[match.id].user.name }}</router-link>
        </td>

      </tr>

    </tbody>
  </table>

  <!--{{ pagination_buttons(matches) }}-->

  </div>
</template>

<script>
export default {
  name: 'matches',
  data () {
    return {
      user: {
        isLoggedIn: false,
        steamid: '',
        userid: ''
      },
      my_matches: true,
      all_matches: false, // TODO
      matches: [],
      matchinfo: {},
      match_owner: { // TODO
        id: 1,
        name: 'hoge'
      },
      teamdatas: {},
      serverdatas: {}
    }
  },
  async created () {
    let self = this
    try {
      this.user = await this.CheckLoggedIn()
    } catch (err) {
      this.user = {
        isLoggedIn: false,
        steamid: '',
        userid: ''
      }
    } finally {
      const matches = await this.GetMatches(this.user.userid)
      for (let i = 0; i < matches.length; i++) {
        this.matchinfo[matches[i].id] = await self.GetMatchData(matches[i].id)
      }
    }
  },
  methods: {
    async GetMatches (userid) {
      return new Promise(async (resolve, reject) => {
        if (!userid) {
          const res = await this.axios.get('/api/v1/GetMatches', {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${this.$store.state.auth}`
            },
            data: null
          })
          this.matches = res.data
          resolve(res.data)
        } else {
          const res = await this.axios.get(`/api/v1/GetMatches?userID=${userid}`, {
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${this.$store.state.auth}`
            },
            data: null
          })
          this.matches = res.data
          resolve(res.data)
        }
      })
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
