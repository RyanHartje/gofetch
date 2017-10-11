{{ define "content" }}
    <div class="container">
      <br>
      <br>
      <br>
      <img class="img-responsive center" src="/assets/images/gold.svg">
      <form class="form-signin">
        <h2 class="form-signin-heading">Open your vault</h2>
        <label for="username" class="sr-only">Username</label>
        <input type="text" id="username" class="form-control" placeholder="Username" required autofocus>
        <label for="password" class="sr-only">Password</label>
        <input type="password" id="password" class="form-control" placeholder="Password" required>
        <div class="checkbox">
          <label>
            <input type="checkbox" value="remember-me"> Remember me
          </label>
        </div>
        <button class="btn btn-lg btn-gold btn-block" type="submit">Log in</button>
      </form>

    </div> <!-- /container -->

{{ end }}
