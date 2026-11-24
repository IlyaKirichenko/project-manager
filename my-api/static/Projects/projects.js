async function loadProjects() {
    const response = await fetch("/project");

    if (!response.ok){
        window.location.href = "/";
        return;
    }

    const projects = await response.json();

    const list = document.getElementById("projectsList");
    list.innerHTML = "";
    if (projects.length === 0){
        list.textContent = "Проектов нет";
        return;
    }
    for (const p of projects){
        const card = document.createElement("div");
        card.className = "project-item";

        card.innerHTML = `
        <div class="project-title">${p.name_project}</div>
        <div>${p.description_project}</div>
        <div class="project-meta">Срок: ${p.deadline}</div>
        <a class="project-link" href="/static/Project/project.html?id=${p.id_project}"> Открыть задачи </a>
        `;
        list.append(card);
    }
}
loadProjects();