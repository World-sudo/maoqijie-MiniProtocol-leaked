package captcha

// captchaPage GeeTest V4 验证码页面 HTML
// 逆向自 sso.mini1.cn: initGeetest4 + captcha.getValidate()
// %s 占位符: captchaID
const captchaPage = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="utf-8">
<title>迷你世界 - 验证码</title>
<style>
body{font-family:system-ui,sans-serif;display:flex;justify-content:center;
align-items:center;min-height:100vh;margin:0;background:#f5f5f5}
.box{background:#fff;border-radius:12px;padding:40px;text-align:center;
box-shadow:0 2px 12px rgba(0,0,0,.1);min-width:320px}
h2{color:#333;margin-bottom:20px}
#status{color:#666;margin-top:16px}
.ok{color:#52c41a!important}
.err{color:#f5222d!important}
</style>
</head>
<body>
<div class="box">
<h2>请完成验证码</h2>
<div id="captcha"></div>
<p id="status">正在加载验证码...</p>
</div>
<script src="https://static.geetest.com/v4/gt4.js"></script>
<script>
initGeetest4({
  captchaId: '%s',
  product: 'bind'
}, function(captcha) {
  document.getElementById('status').textContent = '请点击完成验证';
  captcha.showCaptcha();
  captcha.onSuccess(function() {
    var t = captcha.getValidate();
    document.getElementById('status').textContent = '正在提交...';
    fetch('/callback', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify(t)
    }).then(function(r) {
      if (r.ok) {
        document.getElementById('status').className = 'ok';
        document.getElementById('status').textContent = '验证成功! 可以关闭此页面';
      }
    }).catch(function(e) {
      document.getElementById('status').className = 'err';
      document.getElementById('status').textContent = '提交失败: ' + e;
    });
  });
  captcha.onError(function() {
    document.getElementById('status').className = 'err';
    document.getElementById('status').textContent = '验证码加载失败，请刷新页面';
  });
});
</script>
</body>
</html>`

// successPage 验证成功后的页面
const successPage = `<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>验证成功</title>
<style>body{font-family:system-ui;display:flex;justify-content:center;
align-items:center;min-height:100vh;margin:0;background:#f5f5f5}
.box{background:#fff;border-radius:12px;padding:40px;text-align:center;
box-shadow:0 2px 12px rgba(0,0,0,.1)}
h2{color:#52c41a}</style></head>
<body><div class="box"><h2>验证成功!</h2><p>可以关闭此页面，回到终端查看结果。</p></div></body></html>`
