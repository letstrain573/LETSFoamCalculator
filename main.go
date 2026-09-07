package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

var serverFlag = flag.Bool("server", false, "run the web server")

const appVersion = "v5"

const homePage = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>LETS Foam Calculator</title>
	<style>
		:root {
			--orange: #f28c28;
			--dark: #20252b;
			--light-bg: #f1f3f5;
			--card: #ffffff;
			--text: #222222;
			--muted: #666666;
		}
		* { box-sizing: border-box; }
		body {
			margin: 0;
			font-family: Arial, Helvetica, sans-serif;
			background: var(--light-bg);
			color: var(--text);
		}
		nav {
			background: var(--dark);
			color: #fff;
			padding: 0 24px;
			display: flex;
			align-items: center;
			justify-content: space-between;
			min-height: 64px;
		}
		.brand {
			display: flex;
			align-items: center;
			gap: 12px;
			font-size: 1.25rem;
			font-weight: 700;
		}
		.brand img {
			width: 76px;
			height: auto;
			display: block;
			border-radius: 6px;
			background: #fff;
			padding: 4px;
		}
		.nav-links { display: flex; gap: 22px; }
		.nav-links a {
			color: #fff;
			text-decoration: none;
			font-weight: 600;
			padding: 22px 0 18px;
			border-bottom: 4px solid transparent;
		}
		.nav-links a.active,
		.nav-links a:hover {
			color: var(--orange);
			border-bottom-color: var(--orange);
		}
		main {
			max-width: 1100px;
			margin: 0 auto;
			padding: 34px 20px 48px;
		}
		.logo-card {
			background: var(--card);
			border-radius: 14px;
			padding: 26px;
			box-shadow: 0 4px 16px rgba(0,0,0,0.08);
			text-align: center;
			margin-bottom: 28px;
		}
		.logo-card img {
			max-width: 420px;
			width: 100%;
			height: auto;
			display: block;
			margin: 0 auto;
		}
		.hero { text-align: center; margin-bottom: 28px; }
		.hero h1 { margin-bottom: 12px; font-size: clamp(2rem, 5vw, 3rem); }
		.hero p {
			max-width: 760px;
			margin: 0 auto 22px;
			color: var(--muted);
			font-size: 1.05rem;
			line-height: 1.6;
		}
		.primary-button {
			display: inline-block;
			background: var(--orange);
			color: #fff;
			border-radius: 9px;
			padding: 13px 22px;
			font-size: 1rem;
			font-weight: 700;
			text-decoration: none;
		}
		.feature-grid {
			display: grid;
			grid-template-columns: repeat(3, 1fr);
			gap: 20px;
			margin-top: 30px;
		}
		.feature-card {
			background: var(--card);
			border-radius: 12px;
			padding: 24px;
			box-shadow: 0 3px 12px rgba(0,0,0,0.06);
		}
		.feature-card h2 { margin-top: 0; }
		.feature-card p { color: var(--muted); line-height: 1.5; }
		.app-footer {
			margin-top: 34px;
			text-align: center;
			color: var(--muted);
			font-size: 0.85rem;
		}

		@media (max-width: 760px) {
			nav {
				flex-direction: column;
				align-items: flex-start;
				padding-top: 14px;
			}
			.nav-links { flex-wrap: wrap; }
			.feature-grid { grid-template-columns: 1fr; }
		}
	</style>
</head>
<body>
	<nav>
		<div class="brand">
			<img src="/assets/LOGO.jpg" alt="LETS Training Solutions">
			<span>LETS Foam Calculator</span>
		</div>
		<div class="nav-links">
			<a href="/" class="active">Home</a>
			<a href="/calculator">Calculator</a>
			<a href="#about">About</a>
		</div>
	</nav>

	<main>
		<section class="hero">
			<h1>Welcome to LETS Foam Calculator</h1>
			<p>
				Use this application to calculate foam concentrate and water requirements
				for emergency response planning and operations.
			</p>
			<a class="primary-button" href="/calculator">Open Calculator</a>
		</section>

		<section class="feature-grid">
			<div class="feature-card">
				<h2>Foam Calculator</h2>
				<p>Enter incident information and calculate required foam solution quantities.</p>
			</div>
			<div class="feature-card" id="about">
				<h2>Learn More</h2>
				<p>Use clear inputs and results to support planning, training, and field decision-making.</p>
			</div>
			<div class="feature-card">
				<h2>Be Prepared</h2>
				<p>Use the calculator as a planning aid for emergency response readiness.</p>
			</div>
		</section>

		<footer class="app-footer">Version {{.AppVersion}}</footer>
	</main>
</body>
</html>`

const calculatorPage = `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Calculator | LETS Foam Calculator</title>
	<style>
		:root {
			--orange: #f28c28;
			--dark: #20252b;
			--light-bg: #f1f3f5;
			--card: #ffffff;
			--text: #222222;
			--muted: #666666;
			--border: #cfd4da;
		}
		* { box-sizing: border-box; }
		body {
			margin: 0;
			font-family: Arial, Helvetica, sans-serif;
			background: var(--light-bg);
			color: var(--text);
		}
		nav {
			background: var(--dark);
			color: #fff;
			padding: 0 24px;
			display: flex;
			align-items: center;
			justify-content: space-between;
			min-height: 64px;
		}
		.brand {
			display: flex;
			align-items: center;
			gap: 12px;
			font-size: 1.25rem;
			font-weight: 700;
		}
		.brand img {
			width: 76px;
			height: auto;
			display: block;
			border-radius: 6px;
			background: #fff;
			padding: 4px;
		}
		.nav-links { display: flex; gap: 22px; }
		.nav-links a {
			color: #fff;
			text-decoration: none;
			font-weight: 600;
			padding: 22px 0 18px;
			border-bottom: 4px solid transparent;
		}
		.nav-links a.active,
		.nav-links a:hover {
			color: var(--orange);
			border-bottom-color: var(--orange);
		}
		main {
			max-width: 1100px;
			margin: 0 auto;
			padding: 34px 20px 48px;
		}
		.logo-card,
		.calculator-card {
			background: var(--card);
			border-radius: 14px;
			box-shadow: 0 4px 16px rgba(0,0,0,0.08);
		}
		.logo-card {
			padding: 24px;
			text-align: center;
			margin-bottom: 26px;
		}
		.logo-card img {
			max-width: 420px;
			width: 100%;
			height: auto;
			display: block;
			margin: 0 auto;
		}
		.calculator-logo-card {
			padding: 12px 20px;
		}
		.calculator-logo-card img {
			max-width: 180px;
		}
		.calculator-card {
			padding: 28px;
		}
		h1 {
			margin-top: 0;
			margin-bottom: 8px;
		}
		.intro {
			color: var(--muted);
			margin-top: 0;
			margin-bottom: 26px;
			line-height: 1.5;
		}
		.form-row {
			display: grid;
			grid-template-columns: 1.6fr 1.1fr 1fr 1.1fr;
			gap: 20px;
			align-items: start;
		}
		.form-group {
			display: flex;
			flex-direction: column;
			gap: 8px;
		}
		label {
			font-weight: 700;
			min-height: 24px;
			display: flex;
			align-items: flex-end;
		}
		select,
		input[type="number"] {
			width: 100%;
			min-height: 46px;
			border: 1px solid var(--border);
			border-radius: 8px;
			background: #fff;
			padding: 10px 12px;
			font-size: 1rem;
			color: var(--text);
		}
		select:focus,
		input[type="number"]:focus {
			outline: 2px solid rgba(242, 140, 40, 0.28);
			border-color: var(--orange);
		}
		.field-help {
			font-size: 0.88rem;
			color: var(--muted);
			margin: 0;
		}
		.actions {
			margin-top: 26px;
		}
		.back-link {
			display: inline-block;
			color: var(--orange);
			font-weight: 700;
			text-decoration: none;
		}

		@media (max-width: 760px) {
			nav {
				flex-direction: column;
				align-items: flex-start;
				padding-top: 14px;
			}
			.nav-links { flex-wrap: wrap; }
			.form-row { grid-template-columns: 1fr; }
			.calculator-card { padding: 22px; }
		}
	</style>
</head>
<body>
	<nav>
		<div class="brand">
			<img src="/assets/LOGO.jpg" alt="LETS Training Solutions">
			<span>LETS Foam Calculator</span>
		</div>
		<div class="nav-links">
			<a href="/">Home</a>
			<a href="/calculator" class="active">Calculator</a>
			<a href="/#about">About</a>
		</div>
	</nav>

	<main>
		<section class="calculator-card">
			<h1>Foam Calculator</h1>
			<p class="intro">
				Select the liquid type, incident type, foam concentration, and enter the spill area.
			</p>

			<form>
				<div class="form-row">
					<div class="form-group">
						<label for="liquidType">Flammable and Combustible Liquids</label>
						<select id="liquidType" name="liquidType">
							<option value="hydrocarbon">Hydrocarbon</option>
							<option value="polar-solvent">Polar Solvent</option>
						</select>
					</div>

					<div class="form-group">
						<label for="incidentType">Incident Type</label>
						<select id="incidentType" name="incidentType">
							<option value="spill-fires">Spill Fires</option>
							<option value="in-depth-fires">In Depth Fires</option>
						</select>
					</div>

					<div class="form-group">
						<label for="foamConcentration">Foam Concentration</label>
						<select id="foamConcentration" name="foamConcentration">
							<option value="3">3%</option>
							<option value="6">6%</option>
						</select>
					</div>

					<div class="form-group">
						<label for="spillSize">Spill Size (sq ft)</label>
						<input
							type="number"
							id="spillSize"
							name="spillSize"
							min="100"
							max="100000"
							step="1"
							inputmode="numeric"
							placeholder="100 - 100000"
							required
						>
						<p class="field-help">Enter a whole number from 100 to 100,000 sq ft.</p>
					</div>
				</div>
			</form>

			<div class="actions">
				<a class="back-link" href="/">← Back to Home</a>
			</div>
		</section>
	</main>
</body>
</html>`

func renderHTML(w http.ResponseWriter, page string) {
	data := struct {
		AppVersion string
	}{
		AppVersion: appVersion,
	}

	tmpl, err := template.New("page").Parse(page)
	if err != nil {
		http.Error(w, "Unable to load page", http.StatusInternalServerError)
		log.Printf("template parse error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Printf("template execute error: %v", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	renderHTML(w, homePage)
}

func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/calculator" {
		http.NotFound(w, r)
		return
	}
	renderHTML(w, calculatorPage)
}

func main() {
	flag.Parse()

	if !*serverFlag {
		fmt.Println("Usage:")
		fmt.Println("  lets_foam_calculator.exe -server")
		return
	}

	http.Handle(
		"/assets/",
		http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))),
	)
	http.HandleFunc("/calculator", calculatorHandler)
	http.HandleFunc("/", homeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("LETS Foam Calculator listening on http://localhost:%s", port)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
