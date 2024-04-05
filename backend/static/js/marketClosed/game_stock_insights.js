function animateGameStocks() {
    const showEvent = new Event("show")

    const gameStocks = document.querySelectorAll(".game-stock");

    if (gameStocks.length === 0) {
        throw new Error("No game stocks found")
        return
    }

    gameStocks.forEach((gameStock, index) => {
        let gameStockId = gameStock.getAttribute("data-game-stock-id")

        document.querySelector(`div#game-stock-view-${gameStockId}`).classList.add("first")
    })

    gameStocks.forEach((gameStock, index) => {
        gameStock.addEventListener("show", (e) => {

            document.querySelector("#stock-view").classList.remove("border")
            document.querySelector("#insight-view").classList.add("border")

            let gameStockId = gameStock.getAttribute("data-game-stock-id")

            const gameStockInsights = document.querySelectorAll(`div.game-stock-insight-${gameStockId}`)

            let revealInsightPromises = Array.from(gameStockInsights).map((gameStockInsight, index) => {
                return new Promise((resolve) => {
                    setTimeout(() => {
                        /** @type {NodeListOf<HTMLDivElement>} */
                        const gameInsights = document.querySelectorAll("div.game-insight")
                        gameInsights.forEach(gameInsight => {
                            gameInsight.style.display = "none"
                        })

                        const stockTotalInsightValue = gameStock.querySelector(".stock-total-insight-value")

                        let value = parseFloat(stockTotalInsightValue.innerText)

                        /** @type {HTMLSpanElement?} */
                        const priceModifier = gameStockInsight.querySelector("span.price-modifier")

                        value += parseFloat(priceModifier?.innerText || "0")

                        gameStock.querySelector(".stock-total-insight-value").innerText = value.toFixed(2)

                        gameStockInsights.forEach(gsi => {
                            gsi.querySelector(".stock-total-insight-value").innerText = value.toFixed(2)
                        })


                        gameStockInsight.style.display = "grid"
                        resolve()
                    }, index * 2000);
                })
            })

            Promise.all(revealInsightPromises).then(() => {
                // allow for the final animation to finish running
                return new Promise((resolve) => {
                    setTimeout(() => {
                        const stockTotalInsightValue = gameStock.querySelector(".stock-total-insight-value")

                        console.log("insights complete, rendering stock value animation...");
                        console.log(stockTotalInsightValue)

                        let insightValue = parseFloat(stockTotalInsightValue?.innerText || "0")

                        if (insightValue === 0) {
                            resolve({ insightValue, stockInc: 0 })
                            return
                        }

                        let stockInc = insightValue / Math.abs(insightValue)

                        resolve({ stockInc })
                    }, 1500)
                })
            }).then(({ stockInc }) => {

                return new Promise((resolve) => {
                    document.querySelector("#insight-view").classList.remove("border")
                    document.querySelector("#stock-view").classList.add("border")

                    let intervalId = setInterval(() => {
                        const stockTotalInsightValue = gameStock.querySelector(".stock-total-insight-value")
                        let insightValue = parseFloat(stockTotalInsightValue?.innerText || "0")

                        if (Math.abs(insightValue) <= 0.001) {
                            clearInterval(intervalId)
                            resolve()
                            return
                        }

                        insightValue -= stockInc * 0.1

                        stockTotalInsightValue.innerText = insightValue.toFixed(2)


                        document.querySelectorAll(`span.stock-total-insight-value-${gameStockId}`).forEach(elem => {
                            elem.innerText = insightValue.toFixed(2)
                        })

                        const gameStockValueDisplay = gameStock.querySelector(".stock-value")
                        let gameStockValue = parseFloat(gameStockValueDisplay.innerText)

                        gameStockValueDisplay.innerText = (gameStockValue + stockInc * 0.1).toFixed(2)



                        console.log(stockTotalInsightValue.innerText)
                    }, 50)
                })
            }).then(() => {
                if (index < gameStocks.length - 1) {
                    setTimeout(() => {
                        gameStocks[index + 1].dispatchEvent(showEvent)
                    }, 200)
                }
            })
        })
    })

    gameStocks[0].dispatchEvent(showEvent)
}

// used in hidden game stock section
function animateCurrency() {
    const currencyGameStock = document.querySelector("#game-stock-currency");

    let currencyId = currencyGameStock.getAttribute("data-game-stock-id")

    const gameStockInsights = document.querySelectorAll(`div.game-stock-insight-${currencyId}`)

    let revealInsightPromises = Array.from(gameStockInsights).map((gameStockInsight, index) => {
        return new Promise((resolve) => {
            setTimeout(() => {
                /** @type {NodeListOf<HTMLDivElement>} */
                const gameInsights = document.querySelectorAll("div.game-insight")
                gameInsights.forEach(gameInsight => {
                    gameInsight.style.display = "none"
                })

                const stockTotalInsightValue = currencyGameStock.querySelector(".stock-total-insight-value")

                let value = parseFloat(stockTotalInsightValue.innerText)

                /** @type {HTMLSpanElement?} */
                const priceModifier = gameStockInsight.querySelector("span.price-modifier")

                value += parseFloat(priceModifier?.innerText || "0")

                currencyGameStock.querySelector(".stock-total-insight-value").innerText = value.toFixed(0)

                gameStockInsights.forEach(gsi => {
                    gsi.querySelector(".stock-total-insight-value").innerText = value.toFixed(0)
                })


                gameStockInsight.style.display = "grid"
                resolve()
            }, index * 2000);
        })
    })

    Promise.all(revealInsightPromises).then(() => {
        // allow for the final animation to finish running
        return new Promise((resolve) => {
            setTimeout(() => {
                const currencyTotalInsightValue = currencyGameStock.querySelector(".stock-total-insight-value")

                console.log("insights complete, rendering stock value animation...");
                console.log(currencyTotalInsightValue)

                let cashChangePercent = parseFloat(currencyTotalInsightValue?.innerText || "0")

                resolve({ cashChangePercent })

            }, 1500)
        })
    }).then(({ cashChangePercent }) => {
        document.querySelector("#insight-view").classList.remove("border")
        document.querySelector("#player-view").classList.add("border")

        const players = document.querySelectorAll(".player")

        players.forEach(player => {
            const playerCashChange = player.querySelector(".cash-change")

            const symbol = (cashChangePercent >= 0) ? "+" : "-"

            playerCashChange.innerText = `${symbol}${cashChangePercent}%`
            playerCashChange.classList.add((cashChangePercent >= 0) ? "text-success" : "text-danger")
            playerCashChange.style.display = "block"

            const playerCash = player.querySelector(".player-cash")

            const oldPlayerCashString = playerCash?.innerText || "0"
            const oldPlayerCash = parseFloat(oldPlayerCashString)

            const newPlayerCash = oldPlayerCash * (1 + cashChangePercent / 100)

            playerCash.innerText = newPlayerCash.toFixed(0)

        })
    }).then(() => {
        // wait for the final animation to finish
        setTimeout(() => {
            animateHoldStocks()
        }, 1500)
    })
}

// used in hidden game stock section
function animateHoldStocks() {
    const holdStockGameStock = document.querySelector("#game-stock-hold-stock-price");

    let holdStockId = holdStockGameStock.getAttribute("data-game-stock-id");

    const stockHoldInsights = document.querySelectorAll(`div.game-stock-insight-${holdStockId}`)

    const stockHoldInsight = stockHoldInsights[0];

    const gameInsights = document.querySelectorAll("div.game-insight")
    gameInsights.forEach(gameInsight => {
        gameInsight.style.display = "none"
    })

    stockHoldInsight.style.display = "grid"

    let insightId = stockHoldInsight.getAttribute("data-insight-id")

    const firstStockHoldModal = new bootstrap.Modal(`#stock-hold-modal-${insightId}`, {
        keyboard: false
    })

    firstStockHoldModal.show();

    const gameID = document.querySelector("#gameID").value;

    const socket = new WebSocket(`ws://localhost:4040/hold-stock-waiting/${gameID}`);

    socket.onmessage = function (event) {
        const message = JSON.parse(event.data);
        console.log(message);

        if (!message.hasOwnProperty('nextInsight')) {
            console.error("invalid websocket response")
            return
        }

        const newGameInsight = document.querySelector(`#insight-${message.nextInsight}`);

        if (newGameInsight === null) {
            console.error("invalid game stock insight recieved from websocket")
            return
        }

        gameInsights.forEach(gameInsight => {
            gameInsight.style.display = "none"
        })

        newGameInsight.style.display = "grid";

        const stockHoldModal = new bootstrap.Modal(`#stock-hold-modal-${message.nextInsight}`, {
            keyboard: false
        })

        stockHoldModal.show();


    }
}