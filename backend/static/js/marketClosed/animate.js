const showEvent = new Event("show")

const gameStocks = document.querySelectorAll(".game-stock");

gameStocks.forEach((gameStock, index) => {
    let gameStockId = gameStock.getAttribute("data-game-stock-id")

    document.querySelector(`div#game-stock-view-${gameStockId}`).classList.add("first")
})

gameStocks.forEach((gameStock, index) => {
    gameStock.addEventListener("show", (e) => {

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
                    /* @type {HTMLSpanElement?} */
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
                        resolve({insightValue, stockInc: 0})
                        return
                    }

                    let stockInc = insightValue / Math.abs(insightValue)

                    resolve({stockInc})
                }, 1500)
            })
        }).then(({stockInc}) => {

            return new Promise((resolve) => {

                let intervalId = setInterval(() => {
                    const stockTotalInsightValue = gameStock.querySelector(".stock-total-insight-value")
                    let insightValue = parseFloat(stockTotalInsightValue?.innerText || "0")

                    if(Math.abs(insightValue) <= 0.001){
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
            if(index < gameStocks.length - 1){
                setTimeout(() => {
                    gameStocks[index + 1].dispatchEvent(showEvent)
                }, 200)
            }
        })
    })
})

gameStocks[0].dispatchEvent(showEvent)