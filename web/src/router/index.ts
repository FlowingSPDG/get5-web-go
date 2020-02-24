import Vue from 'vue'
import Router from 'vue-router'
import axios from 'axios'
import VueAxios from 'vue-axios'

import Matches from '@/components/Matches.vue'
import Match from '@/components/Match.vue'
import MatchCreate from '@/components/MatchCreate.vue'
import MyMatches from '@/components/MyMatches.vue'
import MyServers from '@/components/MyServers.vue'
import Metrics from '@/components/Metrics.vue'
import Team from '@/components/Team.vue'
import Teams from '@/components/Teams.vue'
import TeamCreate from '@/components/TeamCreate.vue'
import ServerCreate from '@/components/ServerCreate.vue'
import User from '@/components/User.vue'
import PugCreate from '@/components/PugCreate.vue'

Vue.use(VueAxios, axios)
Vue.use(Router)

export default new Router({
  mode: 'hash',
  routes: [
    {
      path: '/',
      redirect: '/matches'
    },
    {
      path: '/matches',
      name: 'Matches',
      component: Matches
    },
    {
      path: '/matches/:userid',
      name: 'Matches for user',
      component: Matches
    },
    {
      path: '/myservers',
      name: 'My Servers',
      component: MyServers
    },
    {
      path: '/mymatches',
      name: 'My matches',
      component: MyMatches
    },
    {
      path: '/match/create',
      name: 'Create Match',
      component: MatchCreate
    },
    {
      path: '/match/create-pug',
      name: 'Create Pug',
      component: PugCreate
    },
    {
      path: '/match/:matchid',
      name: 'Match',
      component: Match
    },
    {
      path: '/metrics',
      name: 'Metrics',
      component: Metrics
    },
    {
      path: '/team/create',
      name: 'Create Team',
      component: TeamCreate,
      props: {
        edit: false
      }
    },
    {
      path: '/team/:teamid/edit',
      name: 'Edit Team',
      component: TeamCreate,
      props: {
        edit: true
      }
    },
    {
      path: '/team/:teamid',
      name: 'Team',
      component: Team
    },
    {
      path: '/teams/:userid',
      name: 'Teams',
      component: Teams
    },
    {
      path: '/myteams',
      name: 'My Teams',
      component: Teams
    },
    {
      path: '/user/:userid',
      name: 'User',
      component: User
    },
    {
      path: '/server/create',
      name: 'Create Server',
      component: ServerCreate,
      props: {
        edit: false
      }
    },
    {
      path: '/server/:serverID/edit',
      name: 'Edit Server',
      component: ServerCreate,
      props: {
        edit: true
      }
    }
  ]
})
