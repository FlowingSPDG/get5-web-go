<template>
<div id="container">
<div id="content">
  <el-form ref="form" :model="form" :rules="rules" label-position="left" label-width="120px">
  <el-form-item :label="$t('TeamCreate.FormTeamName')" style="width: 653px;" prop="name">
    <el-input v-model="form.name"></el-input>
  </el-form-item>

  <el-form-item :label="$t('TeamCreate.FormTeamTag')" style="width: 653px;" prop="tag">
    <el-input v-model="form.tag"></el-input>
  </el-form-item>

  <el-form-item :label="$t('TeamCreate.FormLogoName')" style="width: 653px;" prop="logo_name">
    <el-select v-model="form.logo_name" placeholder="None" style="width: 653px;">
      <el-option selected value="">None</el-option>
    </el-select>
  </el-form-item>

  <el-form-item :label="$t('TeamCreate.FormCountryFlag')" style="width: 653px;" prop="flag">
    <el-select v-model="form.flag" placeholder="None" style="width: 653px;" filterable>
      <el-option selected value="" >None</el-option>
      <el-option label="Algeria" value="dz">Algeria - DZ</el-option>
      <el-option label="Argentina" value="ar">Argentina - AR</el-option>
      <el-option label="Australia" value="au">Australia - AU</el-option>
      <el-option label="Austria" value="at">Austria - AT</el-option>
      <el-option label="Belarus" value="by">Belarus - BY</el-option>
      <el-option label="Belgium" value="be">Belgium - BE</el-option>
      <el-option label="Brazil" value="br">Brazil - BR</el-option>
      <el-option label="Bulgaria" value="bg">Bulgaria - BG</el-option>
      <el-option label="Canada" value="ca">Canada - CA</el-option>
      <el-option label="Chile" value="cl">Chile - CL</el-option>
      <el-option label="China" value="cn">China - CN</el-option>
      <el-option label="Cocos Islands" value="cc">Cocos Islands - CC</el-option>
      <el-option label="Croatia" value="hr">Croatia - HR</el-option>
      <el-option label="Czech Republic" value="cz">Czech Republic - CZ</el-option>
      <el-option label="Denmark" value="dk">Denmark - DK</el-option>
      <el-option label="Estonia" value="ee">Estonia - EE</el-option>
      <el-option label="European Union" value="eu">European Union - EU</el-option>
      <el-option label="Finland" value="fi">Finland - FI</el-option>
      <el-option label="France" value="fr">France - FR</el-option>
      <el-option label="Germany" value="de">Germany - DE</el-option>
      <el-option label="Greece" value="gr">Greece - GR</el-option>
      <el-option label="Guadeloupe" value="gp">Guadeloupe - GP</el-option>
      <el-option label="Hong Kong" value="hk">Hong Kong - HK</el-option>
      <el-option label="Hungary" value="hu">Hungary - HU</el-option>
      <el-option label="Iceland" value="is">Iceland - IS</el-option>
      <el-option label="India" value="in">India - IN</el-option>
      <el-option label="Indonesia" value="id">Indonesia - ID</el-option>
      <el-option label="Iran" value="ir">Iran - IR</el-option>
      <el-option label="Ireland" value="ie">Ireland - IE</el-option>
      <el-option label="Israel" value="il">Israel - IL</el-option>
      <el-option label="Italy" value="it">Italy - IT</el-option>
      <el-option label="Japan" value="jp">Japan - JP</el-option>
      <el-option label="Kazahkstan" value="kz">Kazahkstan - KZ</el-option>
      <el-option label="Latvia" value="lv">Latvia - LV</el-option>
      <el-option label="Libya" value="ly">Libya - LY</el-option>
      <el-option label="Liechtenstein" value="lt">Liechtenstein - LT</el-option>
      <el-option label="Luxembourg" value="lu">Luxembourg - LU</el-option>
      <el-option label="Macao" value="mo">Macao - MO</el-option>
      <el-option label="Macedonia" value="mk">Macedonia - MK</el-option>
      <el-option label="Malaysia" value="my">Malaysia - MY</el-option>
      <el-option label="Mexico" value="mx">Mexico - MX</el-option>
      <el-option label="Netherlands" value="nl">Netherlands - NL</el-option>
      <el-option label="New Zealand" value="nz">New Zealand - NZ</el-option>
      <el-option label="Norway" value="no">Norway - NO</el-option>
      <el-option label="Pakistan" value="pk">Pakistan - PK</el-option>
      <el-option label="Peru" value="pe">Peru - PE</el-option>
      <el-option label="Phillippines" value="ph">Phillippines - PH</el-option>
      <el-option label="Poland" value="pl">Poland - PL</el-option>
      <el-option label="Portugal" value="pt">Portugal - PT</el-option>
      <el-option label="Reunion" value="re">Reunion - RE</el-option>
      <el-option label="Romania" value="ro">Romania - RO</el-option>
      <el-option label="Russia" value="ru">Russia - RU</el-option>
      <el-option label="Saudi Arabia" value="sa">Saudi Arabia - SA</el-option>
      <el-option label="Serbia" value="rs">Serbia - RS</el-option>
      <el-option label="Singapore" value="sg">Singapore - SG</el-option>
      <el-option label="Slovenia" value="si">Slovenia - SI</el-option>
      <el-option label="Slovokia" value="sk">Slovokia - SK</el-option>
      <el-option label="South Africa" value="za">South Africa - ZA</el-option>
      <el-option label="South Korea" value="kr">South Korea - KR</el-option>
      <el-option label="Spain" value="es">Spain - ES</el-option>
      <el-option label="Sweden" value="se">Sweden - SE</el-option>
      <el-option label="Switzerland" value="ch">Switzerland - CH</el-option>
      <el-option label="Taiwan" value="tw">Taiwan - TW</el-option>
      <el-option label="Thailand" value="th">Thailand - TH</el-option>
      <el-option label="Turkey" value="tr">Turkey - TR</el-option>
      <el-option label="Ukraine" value="ua">Ukraine - UA</el-option>
      <el-option label="United Arab Emirates" value="ae">United Arab Emirates - AE</el-option>
      <el-option label="United Kingdom" value="gb">United Kingdom - GB</el-option>
      <el-option label="United States" value="us">United States - US</el-option>
      <el-option label="Venezuela" value="ve">Venezuela - VE</el-option>
      <el-option label="Vietnam" value="vn">Vietnam - VN</el-option>
    </el-select>
  </el-form-item>

  <el-form-item :label="$t('TeamCreate.FormPlayer')+' 1'" style="width: 653px;" prop="auths[0]">
    <el-input v-model="form.auths[0]"></el-input>
  </el-form-item>
  <el-form-item :label="$t('TeamCreate.FormPlayer')+' 2'" style="width: 653px;" prop="auths[1]">
    <el-input v-model="form.auths[1]"></el-input>
  </el-form-item>
  <el-form-item :label="$t('TeamCreate.FormPlayer')+' 3'" style="width: 653px;" prop="auths[2]">
    <el-input v-model="form.auths[2]"></el-input>
  </el-form-item>
  <el-form-item :label="$t('TeamCreate.FormPlayer')+' 4'" style="width: 653px;" prop="auths[3]">
    <el-input v-model="form.auths[3]"></el-input>
  </el-form-item>
  <el-form-item :label="$t('TeamCreate.FormPlayer')+' 5'" style="width: 653px;" prop="auths[4]">
    <el-input v-model="form.auths[4]"></el-input>
  </el-form-item>
  <el-form-item :label="$t('TeamCreate.FormPlayer')+' 6'" style="width: 653px;" prop="auths[5]">
    <el-input v-model="form.auths[5]"></el-input>
  </el-form-item>
  <el-form-item :label="$t('TeamCreate.FormPlayer')+' 7'" style="width: 653px;" prop="auths[6]">
    <el-input v-model="form.auths[6]"></el-input>
  </el-form-item>

  <el-form-item :label="$t('TeamCreate.FormPublicTeam')" style="width: 653px;">
    <el-switch v-model="form.public_team"></el-switch>
  </el-form-item>

  <el-form-item v-if="!edit" style="width: 653px;">
    <el-button type="primary" @click="RegisterTeam">{{$t('misc.Create')}}</el-button>
  </el-form-item>
  <el-form-item v-if="edit" style="width: 653px;">
    <el-button type="primary" @click="UpdateTeam">{{$t('misc.Update')}}</el-button>
  </el-form-item>
</el-form>

</div>
</div>
</template>

<script>
export default {
  name: 'TeamCreate',
  props: ['edit'],
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
      form: {
        name: '',
        tag: '',
        flag: '',
        logo: '',
        auths: [],
        public_team: false
      },
      rules: {
        name: [
          { required: true, trigger: 'change', min: 5, max: 30, message: 'Length should be 5 to 30' }
        ],
        tag: [
          { required: false, message: 'Please input team tag', trigger: 'change' }
        ],
        flag: [
          { required: false, message: 'Please select country flag', trigger: 'change' }
        ],
        logo: [
          { required: true, message: 'Please pick a logo', trigger: 'change' }
        ],
        auths: [
          { required: true, message: 'Please input your team members steamid', trigger: 'blur' },
          { required: false, message: 'Please input your team members steamid', trigger: 'blur' },
          { required: false, message: 'Please input your team members steamid', trigger: 'blur' },
          { required: false, message: 'Please input your team members steamid', trigger: 'blur' },
          { required: false, message: 'Please input your team members steamid', trigger: 'blur' },
          { required: false, message: 'Please input your team members steamid', trigger: 'blur' },
          { required: false, message: 'Please input your team members steamid', trigger: 'blur' }
        ],
        public_team: [
          { required: false, message: 'Please select team publicity', trigger: 'change' }
        ]
      }
    }
  },
  async created () {
    this.user = await this.CheckLoggedIn()
    if (this.edit) {
      try {
        let res = await this.axios.get(`api/v1/team/${this.$route.params.teamid}/GetTeamInfo`, {
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${this.$store.state.auth}`
          },
          data: null
        })
        this.form.name = res.data.name
        this.form.tag = res.data.tag
        this.form.flag = res.data.flag
        this.form.logo = res.data.logo
        this.form.auths = res.data.steamids
        this.form.public_team = res.data.public_team
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
    async RegisterTeam () {
      const json = JSON.stringify(this.form)
      this.$refs['form'].validate(async (valid) => {
        if (valid) {
          try {
            let res = await this.axios.post('/api/v1/team/create', json)
            this.form = {}
            this.$message({
              message: this.$t('TeamCreate.MessageRegisterSuccess'),
              type: 'success'
            })
            this.$router.push('/myteams')
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
    async UpdateTeam () {
      const json = JSON.stringify(this.form)
      this.$refs['form'].validate(async (valid) => {
        if (valid) {
          try {
            let res = await this.axios.put(`/api/v1/team/${this.$route.params.teamid}/edit`, json)
            this.form = {}
            this.$message({
              message: this.$t('TeamCreate.MessageeEditSuccess'),
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
