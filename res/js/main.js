
function onLoginCode() {
  const email = document.getElementById('email').value
  const state = Math.random().toString(36).slice(2)

  axios({
    method: 'post',
    url: '/api/v1/user/main/code',
    data: {
      email: email,
      state: state,
    },
  })
    .then(function (response) {
      loginToken = response.data
      loginState = state
    })
    .catch(function (error) {
      alert(error.response.data);
    });
}

function loginAuth(callback) {
  const email = document.getElementById('email').value
  const code = document.getElementById('code').value

  axios({
    method: 'post',
    url: '/api/v1/user/main/auth',
    data: {
      email: email,
      state: loginState,
      token: loginToken,
      code: code,
    },
  })
    .then(function (response) {
      var user = response.data
      Cookies.set('access_token', user.accessToken)
      localStorage.setItem("main_token", user.mainToken)

      if (undefined != callback) {
        callback(response)
      } else {
        window.location.href = "/user/oauth" + window.location.search
      }
    })
}

function onAllowAuth() {
  axios({
    method: 'post',
    url: '/api/v1/user/oauth/auth',
    data: {
      mainToken: localStorage.getItem('main_token'),
      clientId: clientId,
      state: state,
    },
  })
    .then(function (response) {
      const code = response.data
      window.location.href = redirect_uri + "?return_to=" + return_to + "&code=" + code
    })
}

function refreshToken(successed, failured) {
  axios({
    method: 'post',
    url: '/api/v1/user/oauth/refresh',
    data: {
      mainToken: localStorage.getItem('main_token'),
    },
  })
    .then(function (response) {
      Cookies.set('access_token', response.data)

      if (undefined != successed) {
        successed(response)
      }
    })
    .catch(function (response) {
      if (undefined != failured) {
        failured(response)
      }
    })
}