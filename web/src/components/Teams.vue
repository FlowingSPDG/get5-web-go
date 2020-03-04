<template>
<div id="content">

  <h1 v-if="my_teams">{{$t('Teams.YourTeams')}}</h1>
  <h1 v-else-if="owner">Teams for <router-link :to="'/user/'+owner.id"> {{ owner.name }}</router-link ></h1>

  <ul class="list-group">
    <li class="list-group-item" v-if="owner.teams.length == 0">
    {{$t('Teams.NoTeamFound')}}
    </li>

    <li class="list-group-item" v-else v-for="(team,index) in owner.teams" :key="index">

      <img :src="get_flag_link(team)" />
      <router-link :to="'/team/'+team.id" class="col-sm-offset-1">{{team.name}}</router-link>

      <div class="pull-right" v-if="CheckTeamDeletable(team)">
        <router-link :to="'/team/'+team.id+'/delete'" class="col-sm-offset-1">{{$t('misc.Delete')}}</router-link>
      </div>

      <div class="pull-right" v-if="CheckTeamEditable(team)">
        <router-link :to="'/team/'+team.id+'/edit'" class="btn btn-primary btn-xs">{{$t('misc.Edit')}}</router-link>
      </div>

    </li>
  </ul>

</div>
</template>

<script>
export default {
  name: 'Teams',
  data () {
    return {
      user: {},
      my_teams: false,
      owner: {
        teams: []
      }
    }
  },
  async created () {
    this.user = await this.CheckLoggedIn()
    if (this.user.userid === this.$route.params.userid || this.$route.path === '/myteams') {
      this.my_teams = true
    }
    if (this.my_teams) {
      this.owner = await this.GetUserData(this.user.userid)
    } else {
      this.owner = await this.GetUserData(this.$route.params.userid)
    }
  },
  methods: {
    CheckTeamEditable: function (team) {
      return team.user_id === this.user.id
    },
    CheckTeamDeletable: function (team) {
      return team.user_id === this.user.id
    },
    async DeleteTeam (teamid) {
      try {
        let res = await this.axios.delete(`/api/v1/team/${teamid}/delete`)
        this.$message({
          message: 'Successfully deleted team.',
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

.fade-enter-active, .fade-leave-active {
  transition: opacity .5s
}
.fade-enter, .fade-leave {
  opacity: 0
}
</style>
