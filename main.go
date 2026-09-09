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

const appVersion = "v6"

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
				Use this application to calculate foam solution, foam concentrate, and water
				requirements for Class B firefighting training and emergency response planning.
			</p>
			<a class="primary-button" href="/calculator">Open Calculator</a>
		</section>

		<section class="feature-grid">
			<div class="feature-card">
				<h2>Foam Calculator</h2>
				<p>Enter incident information and calculate required foam solution quantities.</p>
			</div>
			<div class="feature-card" id="about">
				<h2>Dynamic Defaults</h2>
				<p>Application rate and design duration are preloaded from the selected liquid and incident types.</p>
			</div>
			<div class="feature-card">
				<h2>Advanced Options</h2>
				<p>Experienced users can review and manually override the preloaded application parameters when appropriate.</p>
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
			--orange-dark: #d97817;
			--dark: #20252b;
			--light-bg: #f1f3f5;
			--card: #ffffff;
			--text: #222222;
			--muted: #666666;
			--border: #cfd4da;
			--result-bg: #f8f9fa;
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
		.calculator-card {
			background: var(--card);
			border-radius: 14px;
			box-shadow: 0 4px 16px rgba(0,0,0,0.08);
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
			line-height: 1.4;
		}
		.advanced-card {
			margin-top: 24px;
			border: 1px solid var(--border);
			border-radius: 10px;
			overflow: hidden;
			background: #fff;
		}
		.advanced-card summary {
			cursor: pointer;
			list-style: none;
			padding: 16px 18px;
			font-weight: 700;
			background: #f8f9fa;
			display: flex;
			align-items: center;
			justify-content: space-between;
			gap: 16px;
		}
		.advanced-card summary::-webkit-details-marker { display: none; }
		.advanced-card summary::after {
			content: "+";
			font-size: 1.35rem;
			color: var(--orange);
		}
		.advanced-card[open] summary::after { content: "−"; }
		.advanced-content {
			padding: 18px;
			border-top: 1px solid var(--border);
		}
		.advanced-intro {
			color: var(--muted);
			margin: 0 0 16px;
			line-height: 1.5;
		}
		.advanced-grid {
			display: grid;
			grid-template-columns: repeat(2, minmax(0, 1fr));
			gap: 20px;
		}
		.button-row {
			display: flex;
			align-items: center;
			gap: 18px;
			flex-wrap: wrap;
			margin-top: 26px;
		}
		.calculate-button {
			border: 0;
			background: var(--orange);
			color: #fff;
			border-radius: 9px;
			padding: 13px 24px;
			font-size: 1rem;
			font-weight: 700;
			cursor: pointer;
		}
		.calculate-button:hover { background: var(--orange-dark); }
		.back-link {
			color: var(--orange);
			font-weight: 700;
			text-decoration: none;
		}
		.results-card {
			display: none;
			margin-top: 26px;
			padding: 22px;
			border-radius: 10px;
			background: var(--result-bg);
			border: 1px solid var(--border);
		}
		.results-card.visible { display: block; }
		.results-card h2 { margin: 0 0 18px; }
		.results-grid {
			display: grid;
			grid-template-columns: repeat(3, minmax(0, 1fr));
			gap: 16px;
		}
		.result-item {
			background: #fff;
			border: 1px solid var(--border);
			border-radius: 9px;
			padding: 18px;
		}
		.result-label {
			color: var(--muted);
			font-size: 0.9rem;
			margin-bottom: 8px;
		}
		.result-value {
			font-size: 1.35rem;
			font-weight: 700;
		}
		.result-basis {
			margin: 18px 0 0;
			color: var(--muted);
			font-size: 0.9rem;
			line-height: 1.5;
		}
		.error-message {
			display: none;
			margin-top: 18px;
			padding: 12px 14px;
			border-radius: 8px;
			background: #fff3f3;
			border: 1px solid #d99;
			color: #8a1f1f;
			line-height: 1.4;
		}
		.error-message.visible { display: block; }
		.app-footer {
			margin-top: 26px;
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
			.form-row,
			.advanced-grid,
			.results-grid { grid-template-columns: 1fr; }
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
				Standard application rate and design duration values are loaded automatically.
			</p>

			<form id="foamCalculatorForm">
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

				<details class="advanced-card" id="advancedOptions">
					<summary>Advanced Options</summary>
					<div class="advanced-content">
						<p class="advanced-intro">
							These values are automatically preloaded from the selected liquid and incident types.
							Advanced users may manually override either value before calculating.
						</p>

						<div class="advanced-grid">
							<div class="form-group">
								<label for="applicationRate">Application Rate (GPM/sq ft)</label>
								<input
									type="number"
									id="applicationRate"
									name="applicationRate"
									min="0.01"
									step="0.01"
									required
								>
								<p class="field-help">Automatically loaded; editable when a different approved rate is required.</p>
							</div>

							<div class="form-group">
								<label for="designDuration">Design Duration (minutes)</label>
								<input
									type="number"
									id="designDuration"
									name="designDuration"
									min="1"
									step="1"
									required
								>
								<p class="field-help">Automatically loaded; editable when a different approved duration is required.</p>
							</div>
						</div>
					</div>
				</details>

				<div id="errorMessage" class="error-message" role="alert"></div>

				<div class="button-row">
					<button type="submit" class="calculate-button">Calculate</button>
					<a class="back-link" href="/">← Back to Home</a>
				</div>
			</form>

			<section id="resultsCard" class="results-card" aria-live="polite">
				<h2>Calculation Results</h2>

				<div class="results-grid">
					<div class="result-item">
						<div class="result-label">Total Foam Solution</div>
						<div class="result-value"><span id="totalFoamSolution">0</span> gal</div>
					</div>

					<div class="result-item">
						<div class="result-label">Foam Concentrate Needed</div>
						<div class="result-value"><span id="foamConcentrate">0</span> gal</div>
					</div>

					<div class="result-item">
						<div class="result-label">Water Needed</div>
						<div class="result-value"><span id="waterNeeded">0</span> gal</div>
					</div>
				</div>

				<p id="resultBasis" class="result-basis"></p>
			</section>

			<footer class="app-footer">Version {{.AppVersion}}</footer>
		</section>
	</main>

	<script>
		(function () {
			const defaults = {
				"hydrocarbon|spill-fires": { applicationRate: 0.10, designDuration: 15 },
				"hydrocarbon|in-depth-fires": { applicationRate: 0.10, designDuration: 30 },
				"polar-solvent|spill-fires": { applicationRate: 0.10, designDuration: 15 },
				"polar-solvent|in-depth-fires": { applicationRate: 0.20, designDuration: 60 }
			};

			const form = document.getElementById("foamCalculatorForm");
			const liquidType = document.getElementById("liquidType");
			const incidentType = document.getElementById("incidentType");
			const foamConcentration = document.getElementById("foamConcentration");
			const spillSize = document.getElementById("spillSize");
			const applicationRate = document.getElementById("applicationRate");
			const designDuration = document.getElementById("designDuration");
			const resultsCard = document.getElementById("resultsCard");
			const errorMessage = document.getElementById("errorMessage");

			function loadDefaults() {
				const key = liquidType.value + "|" + incidentType.value;
				const selectedDefaults = defaults[key];

				if (!selectedDefaults) {
					return;
				}

				applicationRate.value = selectedDefaults.applicationRate.toFixed(2);
				designDuration.value = selectedDefaults.designDuration;
				resultsCard.classList.remove("visible");
				errorMessage.classList.remove("visible");
			}

			function formatGallons(value) {
				return value.toLocaleString("en-US", {
					minimumFractionDigits: 2,
					maximumFractionDigits: 2
				});
			}

			function showError(message) {
				errorMessage.textContent = message;
				errorMessage.classList.add("visible");
				resultsCard.classList.remove("visible");
			}

			liquidType.addEventListener("change", loadDefaults);
			incidentType.addEventListener("change", loadDefaults);

			form.addEventListener("submit", function (event) {
				event.preventDefault();

				const area = Number(spillSize.value);
				const rate = Number(applicationRate.value);
				const duration = Number(designDuration.value);
				const concentrationPercent = Number(foamConcentration.value);
				const concentrationFraction = concentrationPercent / 100;

				if (!Number.isFinite(area) || area < 100 || area > 100000) {
					showError("Enter a spill size from 100 to 100,000 square feet.");
					spillSize.focus();
					return;
				}

				if (!Number.isFinite(rate) || rate <= 0) {
					showError("Application Rate must be greater than zero.");
					applicationRate.focus();
					return;
				}

				if (!Number.isFinite(duration) || duration <= 0) {
					showError("Design Duration must be greater than zero.");
					designDuration.focus();
					return;
				}

				if (!Number.isFinite(concentrationPercent) || concentrationPercent <= 0 || concentrationPercent >= 100) {
					showError("Select a valid foam concentration.");
					foamConcentration.focus();
					return;
				}

				const totalFoamSolution = area * rate * duration;
				const foamConcentrateNeeded = totalFoamSolution * concentrationFraction;
				const waterNeeded = totalFoamSolution - foamConcentrateNeeded;

				document.getElementById("totalFoamSolution").textContent = formatGallons(totalFoamSolution);
				document.getElementById("foamConcentrate").textContent = formatGallons(foamConcentrateNeeded);
				document.getElementById("waterNeeded").textContent = formatGallons(waterNeeded);

				document.getElementById("resultBasis").textContent =
					"Calculated using " +
					area.toLocaleString("en-US") + " sq ft × " +
					rate.toFixed(2) + " GPM/sq ft × " +
					duration.toLocaleString("en-US") + " minutes at " +
					concentrationPercent + "% foam concentration.";

				errorMessage.classList.remove("visible");
				resultsCard.classList.add("visible");
				resultsCard.scrollIntoView({ behavior: "smooth", block: "nearest" });
			});

			loadDefaults();
		})();
	</script>
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
	log.Printf("LETS Foam Calculator %s listening on http://localhost:%s", appVersion, port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
