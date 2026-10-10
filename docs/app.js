const catalogElement = document.querySelector("#catalog");
const catalogCount = document.querySelector("#catalog-count");
const resultsElement = document.querySelector("#results");
const resultCount = document.querySelector("#result-count");
const idsInput = document.querySelector("#item-ids");
const requestStatus = document.querySelector("#request-status");
const requestButton = document.querySelector("#send-request");
const apiSettings = document.querySelector("#api-settings");
const apiURLInput = document.querySelector("#api-url");
const sourceInputs = document.querySelectorAll('input[name="source"]');

let sampleItems = [];

apiURLInput.value = localStorage.getItem("items-api-base-url") || "";

for (const input of sourceInputs) {
  input.addEventListener("change", () => {
    const usingAPI = getSource() === "api";
    apiSettings.hidden = !usingAPI;
    setStatus(
      usingAPI
        ? "Live mode sends requests to your configured API."
        : "Sample data is ready to explore.",
    );
  });
}

apiURLInput.addEventListener("change", () => {
  localStorage.setItem("items-api-base-url", apiURLInput.value.trim());
});

requestButton.addEventListener("click", sendRequest);

function getSource() {
  return document.querySelector('input[name="source"]:checked').value;
}

function setStatus(message, kind = "") {
  requestStatus.textContent = message;
  requestStatus.dataset.kind = kind;
}

function parseIDs(value) {
  return [...new Set(value.split(",").map((id) => id.trim()).filter(Boolean))];
}

function formatPrice(price) {
  return new Intl.NumberFormat(undefined, {
    style: "currency",
    currency: "USD",
  }).format(price);
}

function makeElement(tag, className, text) {
  const element = document.createElement(tag);
  if (className) element.className = className;
  if (text !== undefined) element.textContent = text;
  return element;
}

function makeItemCard(item) {
  const card = makeElement("article", "item-card");
  const art = makeElement("div", "item-art", item.name.slice(0, 2).toUpperCase());
  art.setAttribute("aria-hidden", "true");
  card.append(art);

  const body = makeElement("div", "item-card-body");
  const meta = makeElement("div", "item-meta");
  meta.append(
    makeElement("span", "item-id", item.id),
    makeElement("span", "rating", `★ ${item.rating.toFixed(1)}`),
  );
  body.append(meta, makeElement("h3", "", item.name));
  body.append(makeElement("p", "", item.description));
  body.append(makeElement("span", "item-price", formatPrice(item.price)));

  if (item.specifications.length > 0) {
    const specifications = makeElement("ul", "spec-list");
    for (const specification of item.specifications) {
      specifications.append(makeElement("li", "", specification));
    }
    body.append(specifications);
  }

  card.append(body);
  return card;
}

function showResults(items) {
  resultsElement.replaceChildren();
  resultCount.textContent = String(items.length);

  if (items.length === 0) {
    const empty = makeElement("div", "empty-state");
    empty.append(
      makeElement("span", "empty-icon", "{ }"),
      makeElement("p", "", "No matching items were returned."),
    );
    resultsElement.append(empty);
    return;
  }

  for (const item of items) {
    resultsElement.append(makeItemCard(item));
  }
}

function makeCatalogCard(item) {
  const card = makeElement("article", "catalog-card");
  const art = makeElement("div", "item-art", item.name.slice(0, 2).toUpperCase());
  art.setAttribute("aria-hidden", "true");
  card.append(art);

  const body = makeElement("div", "catalog-card-body");
  const meta = makeElement("div", "catalog-meta");
  meta.append(
    makeElement("span", "catalog-id", item.id),
    makeElement("span", "rating", `★ ${item.rating.toFixed(1)}`),
  );
  body.append(meta, makeElement("h3", "", item.name));
  body.append(makeElement("p", "", item.description));
  body.append(makeElement("span", "item-price", formatPrice(item.price)));

  const addButton = makeElement("button", "add-id-button", "+ Add ID to request");
  addButton.type = "button";
  addButton.addEventListener("click", () => {
    const ids = parseIDs(idsInput.value);
    if (!ids.includes(item.id)) ids.push(item.id);
    idsInput.value = ids.join(", ");
    idsInput.focus();
    setStatus(`Added ${item.id} to your request.`, "success");
  });
  body.append(addButton);
  card.append(body);
  return card;
}

async function loadCatalog() {
  try {
    const response = await fetch("./items.json");
    if (!response.ok) {
      throw new Error(`Could not load sample catalog (HTTP ${response.status}).`);
    }

    sampleItems = await response.json();
    if (!Array.isArray(sampleItems)) {
      throw new Error("The sample catalog has an invalid format.");
    }

    catalogCount.textContent = String(sampleItems.length);
    catalogElement.replaceChildren(...sampleItems.map(makeCatalogCard));
  } catch (error) {
    catalogCount.textContent = "!";
    catalogElement.replaceChildren(
      makeElement("p", "catalog-message", error.message),
    );
    setStatus(error.message, "error");
  }
}

async function sendRequest() {
  const ids = parseIDs(idsInput.value);
  if (ids.length === 0) {
    setStatus("Enter at least one item ID to send a request.", "error");
    idsInput.focus();
    return;
  }

  requestButton.disabled = true;
  setStatus("Sending request…");

  try {
    let items;
    if (getSource() === "demo") {
      items = sampleItems.filter((item) => ids.includes(item.id));
      if (items.length === 0) {
        setStatus("No sample items matched those IDs.", "error");
      } else {
        setStatus(`Loaded ${items.length} item${items.length === 1 ? "" : "s"} from sample data.`, "success");
      }
    } else {
      const baseURL = apiURLInput.value.trim();
      if (!baseURL) {
        throw new Error("Enter the base URL of your running API.");
      }

      let endpoint;
      try {
        endpoint = new URL("/v1/api/items", baseURL);
        if (!["http:", "https:"].includes(endpoint.protocol)) {
          throw new Error("The API URL must use HTTP or HTTPS.");
        }
      } catch (error) {
        throw new Error(error.message || "Enter a valid API base URL.");
      }

      localStorage.setItem("items-api-base-url", baseURL);
      const response = await fetch(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ data: ids.map((id) => ({ id })) }),
      });

      let payload;
      try {
        payload = await response.json();
      } catch {
        throw new Error(`The API returned an unreadable response (HTTP ${response.status}).`);
      }

      if (!response.ok) {
        throw new Error(payload.error || `The API returned HTTP ${response.status}.`);
      }
      if (!Array.isArray(payload.data)) {
        throw new Error("The API response is missing its data array.");
      }

      items = payload.data;
      setStatus(`Received ${items.length} item${items.length === 1 ? "" : "s"} from the live API.`, "success");
    }

    showResults(items);
  } catch (error) {
    showResults([]);
    setStatus(error.message, "error");
  } finally {
    requestButton.disabled = false;
  }
}

loadCatalog();
