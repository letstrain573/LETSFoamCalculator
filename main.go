package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

const page = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>LETS Foam Calculator</title>
<style>
*{box-sizing:border-box}
body{margin:0;font-family:Arial,Helvetica,sans-serif;background:#f4f7f9;color:#18202a}
nav{background:#202830;color:white;min-height:76px;padding:0 42px;display:flex;align-items:center;justify-content:space-between;gap:24px}
.brand{font-size:30px;font-weight:700}
.nav-links{display:flex;gap:20px;align-items:center}
.nav-links a{color:white;text-decoration:none;font-size:18px;padding:14px 22px;border-radius:8px}
.nav-links a.active{background:#ff7518}
main{width:min(1400px,92%);margin:30px auto 60px}
.logo-card{background:white;border-radius:14px;min-height:360px;display:flex;align-items:center;justify-content:center;padding:24px;box-shadow:0 5px 20px rgba(0,0,0,.08)}
.logo-card img{display:block;max-width:280px;width:100%;height:auto;object-fit:contain}
.hero{text-align:center;padding:30px 20px 20px}
.hero h1{font-size:clamp(32px,4vw,54px);margin:0 0 12px}
.hero p{max-width:900px;margin:0 auto;font-size:22px;line-height:1.5;color:#687484}
.cards{display:grid;grid-template-columns:repeat(3,1fr);gap:24px;margin-top:10px}
.card{background:white;border-radius:14px;padding:26px;text-align:center;box-shadow:0 5px 20px rgba(0,0,0,.08)}
.icon{width:76px;height:76px;margin:0 auto 12px;border-radius:50%;background:#ff7518;color:white;display:flex;align-items:center;justify-content:center;font-size:34px;font-weight:bold}
.card h2{margin:8px 0;font-size:26px}
.card p{color:#687484;font-size:18px;line-height:1.45;min-height:78px}
.button{display:block;width:100%;padding:14px 18px;margin-top:18px;border-radius:8px;border:2px solid #ff7518;background:white;color:#ff7518;font-size:18px;font-weight:700;text-decoration:none}
.button.primary{background:#ff7518;color:white}
@media(max-width:800px){nav{padding:18px 22px;flex-direction:column}.brand{font-size:25px}.nav-links{flex-wrap:wrap;justify-content:center}.logo-card{min-height:260px}.logo-card img{max-width:210px}.cards{grid-template-columns:1fr}}
</style>
</head>
<body>
<nav>
<div class="brand">LETS Foam Calculator</div>
<div class="nav-links"><a class="active" href="/">Home</a><a href="#calculator">Calculator</a><a href="#about">About</a></div>
</nav>
<main>
<section class="logo-card"><img src="/assets/LOGO.jpg" alt="LETS Training Solutions logo"></section>
<section class="hero">
<h1>Welcome to LETS Foam Calculator</h1>
<p>A simple, reliable tool to help calculate foam concentrate and water requirements for effective emergency response.</p>
</section>
<section class="cards">
<div class="card" id="calculator"><div class="icon">▦</div><h2>Foam Calculator</h2><p>Calculate foam concentrate and water requirements for a variety of scenarios.</p><a class="button primary" href="#calculator">Open Calculator</a></div>
<div class="card" id="about"><div class="icon">▤</div><h2>Learn More</h2><p>Explore foam basics, best practices, and helpful reference information.</p><a class="button" href="#about">About</a></div>
<div class="card"><div class="icon">◆</div><h2>Be Prepared</h2><p>Support safer operations through training, planning, and readiness.</p><a class="button" href="#calculator">Get Started</a></div>
</section>
</main>
</body>
</html>`

func main() {
	server := flag.Bool("server", false, "Start the web server")
	flag.Parse()
	if !*server {
		fmt.Println("LETS Foam Calculator")
		fmt.Println("Use -server to start the web service.")
		fmt.Println("Example: ./lets_foam_calculator.exe -server")
		return
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, page)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port
	log.Printf("LETS Foam Calculator listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
