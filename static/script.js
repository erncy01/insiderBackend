// script.js handles frontend interactions with the backend API.
async function simulateWeek() {
  const res = await fetch("/simulate-week");
  const data = await res.json();
  document.getElementById("log").innerText = JSON.stringify(data, null, 2);
  updateStandings();
  updateMatches();
}

async function simulateAll() {
  const res = await fetch("/simulate-all");
  const data = await res.json();
  document.getElementById("log").innerText = JSON.stringify(data, null, 2);
  updateStandings();
  updateMatches();
}

async function resetLeague() {
  await fetch("/reset");
  document.getElementById("log").innerText = "League has been reset.";
  updateStandings();
  updateMatches();
}

async function editMatch(event) {
  event.preventDefault();

  const body = {
    week: parseInt(document.getElementById("editWeek").value),
    home_team: document.getElementById("editHomeTeam").value,
    away_team: document.getElementById("editAwayTeam").value,
    home_goals: parseInt(document.getElementById("editHomeGoals").value),
    away_goals: parseInt(document.getElementById("editAwayGoals").value)
  };

  const res = await fetch("/edit-match", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });

  const text = await res.text();
  document.getElementById("log").innerText = text;

  updateStandings();
  updateMatches();
}


async function updateStandings() {
  const res = await fetch("/standings");
  const data = await res.json();
  document.getElementById("standings").innerText = JSON.stringify(data, null, 2);
}

async function updateMatches() {
  const res = await fetch("/matches");
  const data = await res.json();
  document.getElementById("matches").innerText = JSON.stringify(data, null, 2);
}

// Load initial data
updateStandings();
updateMatches();
