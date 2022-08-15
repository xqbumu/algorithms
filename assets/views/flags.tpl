<!DOCTYPE html>
<html>

<head>
  <meta http-equiv="content-type" content="text/html; charset=UTF-8">
  <style>
    img {
      display: inline-block;
      margin: -4px -1px;
      padding: 0;
    }

    div#main {
      font-size: 0;
      line-height: 0;
      position: relative;
    }

    div#timer-wrapper {
      position: absolute;
      top: 112px;
      left: 44px;
      width: 198px;
      height: 48px;
    }

    div#timer-background {
      width: 100%;
      height: 100%;
      background-color: rgba(0, 0, 0, 0.5);
      text-align: center;
    }

    span#timer {
      text-align: center;
      font-size: 36px;
      line-height: 48px;
      color: #FFF;
      font-weight: bold;
    }
  </style>
</head>

<body>
  <script type="text/javascript">
    /* <![CDATA[ */
    var start = new Date().getTime();
/* ]]> */
  </script>
  <div id="main">
  {{- range $index, $value := .flags }}
    <img src="{{$value}}" width="24px" height="24px" border="0">
    {{- if br $index 16}}<br>{{end -}}
  {{end}}
    <div id="timer-wrapper"></div>
  </div>
  <script type="text/javascript">
    /* <![CDATA[ */
    window.onload = function () {
      var end = new Date().getTime();
      var element = document.getElementById("timer-wrapper");
      element.innerHTML = '<div id="timer-background"><span id="timer">' + (end - start) + 'ms</span></div>';
    }
/* ]]> */
  </script>

</body>

</html>