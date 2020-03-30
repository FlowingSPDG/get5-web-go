<template>
<div id="app">
    <div id="header">
        <nav class="navbar navbar-default">
            <div class="container-fluid">
                <div class="navbar-header">
                    <button type="button" class="navbar-toggle" data-toggle="collapse" data-target="#myNavbar">
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                        <span class="icon-bar"></span>
                    </button>
                    <transition name="fade">
                        <router-link v-if="LogoTransition" to="/" class="navbar-brand" >{{ $t("App.title") }}</router-link>
                    </transition>
                </div>
                <div class="collapse navbar-collapse" id="myNavbar">
                    <el-menu
                        class="nav navbar-nav"
                        :default-active="activeIndex"
                        mode="horizontal"
                        router
                    >
                        <el-menu-item index="Matches" id="matches" :route="{ path:'/matches' }" v-if="user.isLoggedIn">{{ $t("App.AllMatches") }}</el-menu-item>
                        <el-menu-item index="My matches" id="mymatches" :route="{ path:'/mymatches' }" v-if="user.isLoggedIn">{{ $t("App.MyMatches") }}</el-menu-item>
                        <el-menu-item index="match_create" id="match_create" :route="{ path:'/match/create' }" v-if="user.isLoggedIn">{{ $t("App.CreateMatch") }}</el-menu-item>
                        <el-menu-item index="myteams" id="myteams" :route="{ path:'/myteams' }" v-if="user.isLoggedIn">{{ $t("App.MyTeams") }}</el-menu-item>
                        <el-menu-item index="team_create" id="matches" :route="{ path:'/team/create' }" v-if="user.isLoggedIn">{{ $t("App.CreateTeam") }}</el-menu-item>
                        <el-menu-item index="myservers" id="myservers" :route="{ path:'/myservers' }" v-if="user.isLoggedIn">{{ $t("App.MyServers") }}</el-menu-item>
                        <el-menu-item index="server_create" id="server_create" :route="{ path:'/server/create' }" v-if="user.isLoggedIn">{{ $t("App.AddServer") }}</el-menu-item>
                        <el-menu-item index="logout" id="logout" v-if="user.isLoggedIn"> <a href='/auth/logout' v-if="user.isLoggedIn">{{ $t("App.Logout") }}</a> </el-menu-item>
                        <el-menu-item index="login" id="login" v-if="!user.isLoggedIn"> <a href='/auth/login' v-if="!user.isLoggedIn"> <img src="/img/login_small.png" height="18" /></a> </el-menu-item>
                      <el-dropdown @command="handleLanguage">
                        <el-button>
                          {{ $t("lang.ChangeLanguage") }}<i class="el-icon-arrow-down el-icon--right"></i>
                        </el-button>
                        <el-dropdown-menu slot="dropdown">
                          <el-dropdown-item command="en">English</el-dropdown-item>
                          <el-dropdown-item command="ja">Japanese</el-dropdown-item>
                          <el-dropdown-item command="da">Danish</el-dropdown-item>
                        </el-dropdown-menu>
                      </el-dropdown>
                    </el-menu>
                </div>
            </div>
        </nav>
    </div>
    <router-view />
    <div class="panel-footer text-muted">
        <div>
            Powered by <a href="http://steampowered.com">Steam</a> -
            <router-link id="metrics" to="/metrics" >Stats</router-link>
            <div v-if="version">- Version <a href="https://github.com/FlowingSPDG/get5-web-go">{{ version }}</a></div>
        </div>
    </div>
</div>
</template>

<script>
import axios from 'axios'
export default {
  name: 'App',
  data () {
    return {
      LogoTransition: false,
      version: '',
      activeIndex: '',
      user: {
        isLoggedIn: false,
        adminaccess: false,
        steamid: '',
        userid: ''
      } // should be object from JSON response
    }
  },
  async mounted () {
    this.LogoTransition = true
    this.activeIndex = this.$route.name
    this.user = await this.CheckLoggedIn()
    let Version = await axios.get('/api/v1/GetVersion', {
      headers: this.GenerateHeaders(),
      data: null
    })
    this.version = Version.data.version
    // console.log(this.$store.state.auth)
  },
  methods: {
    handleLanguage: function (command) {
      this.ChangeLanguage(command)
    }
  }
}
</script>

<style>
#app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin-top: 60px;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 5s;
}
.fade-enter, .fade-leave-to /* .fade-leave-active below version 2.1.8 */ {
  opacity: 0;
}

.container {
  padding-right: 15px;
  padding-left: 15px;
  margin-right: auto;
  margin-left: auto;
}
</style>
