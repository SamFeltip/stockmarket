//@ts-check

const html = String.raw

export class StockCalculator extends HTMLElement {

    static get observedAttributes() {
        return [
            "playerStock-quantity", //={fmt.Sprintf("%d", player_stock.Quantity)}
            "gameStock-value",
            "gameStock-sharesAvailable", //={strconv.Itoa(player_stock.GameStock.SharesAvailable())}
            "player-cash", //={strconv.Itoa(player_stock.Player.Cash)}
            "playerStock-quantity-add"
        ]
    }

    attributeChangedCallback(name, oldValue, newValue) {
        console.log(`Attribute ${name} changed from ${oldValue} to ${newValue}`);
    }

    get gameStockValue() {
        let game_stock_value_attr = this.getAttribute('gamestock-value')

        if (game_stock_value_attr == null) {
            throw new Error("gamestock-value attribute is required")
        }

        return parseFloat(game_stock_value_attr)
    }

    get addPlayerStockQuantity() {
        let player_stock_quantity_add_attr = this.getAttribute('playerstock-quantity-add')

        if (player_stock_quantity_add_attr == null) {
            throw new Error("playerstock-quantity-add attribute is required")
        }

        return parseInt(player_stock_quantity_add_attr)
    }

    set addPlayerStockQuantity(value) {

        let rounded_stock_count = Math.floor(value / 1000) * 1000

        /** @type {HTMLInputElement?} */
        let input_stock_count = this.querySelector('input.input-stock')
        if (input_stock_count == null) {
            throw new Error("input-stock element is required")
        }

        input_stock_count.value = rounded_stock_count.toString()

        /** @type {HTMLInputElement?} */
        let input_total_cost = this.querySelector('input.input-cash')

        if (input_total_cost == null) {
            throw new Error("input-cash element is required")
        }

        input_total_cost.value = (rounded_stock_count * this.gameStockValue).toString()

        this.setAttribute('playerStock-quantity-add', rounded_stock_count.toString())

        /* @type {HTMLDivElement?} */
        let transaction_body_elem = this.querySelector('.transaction-body')

        if (transaction_body_elem == null) {
            throw new Error("transaction-body element is required")
        }

        transaction_body_elem.outerHTML = this.transaction_body
    }

    set totalCost(value) {

        let rounded_to_value_stock_count = Math.floor(value / this.gameStockValue);
        let new_stock_count = Math.floor(rounded_to_value_stock_count / 1000) * 1000
        this.addPlayerStockQuantity = new_stock_count

    }
    get totalCost() {
        return this.addPlayerStockQuantity * this.gameStockValue
    }

    get playerStockQuantity() {
        let player_stock_quantity_attr = this.getAttribute('playerstock-quantity')

        if (player_stock_quantity_attr == null) {
            throw new Error("playerstock-quantity attribute is required")
        }

        return parseInt(player_stock_quantity_attr)
    }

    get gameStockSharesAvailable() {
        let game_stock_shares_available_attr = this.getAttribute('gamestock-sharesavailable')

        if (game_stock_shares_available_attr == null) {
            throw new Error("gamestock-sharesavailable attribute is required")
        }

        return parseInt(game_stock_shares_available_attr)
    }

    get playerCash() {
        let player_cash_attr = this.getAttribute('player-cash')

        if (player_cash_attr == null) {
            throw new Error("player-cash attribute is required")
        }

        return parseInt(player_cash_attr)
    }

    constructor() {
        super();
        if (this.getAttribute('playerstock-quantity') == null) {
            throw new Error("playerstock-quantity attribute is required")
        }

        if (this.getAttribute('gamestock-value') == null) {
            throw new Error("gamestock-value attribute is required")
        }

        if (this.getAttribute('gamestock-sharesavailable') == null) {
            throw new Error("gamestock-sharesavailable attribute is required")
        }

        if (this.getAttribute('player-cash') == null) {
            throw new Error("player-cash attribute is required")
        }

        this.addPlayerStockQuantity = 0;
    }

    get transaction_body() {
        console.log("rerendering transaction body");
        return html`
            <div class="transaction-body" style="display: flex; flex-direction: row;">
                <div style="flex: 1;">
                    <h6>
                        Current
                    </h6>
                    
                    <p>
                        ${this.playerStockQuantity} Shares = £${this.playerStockQuantity * this.gameStockValue}
                    </p>
                    <p>
                        ${this.gameStockSharesAvailable} Stocks left to buy
                    </p>
                    <p>
                        £${this.playerCash} Cash available
                    </p>
                </div>

                <div style="flex: 1;">
                    <h6>
                        After transaction
                    </h6>
                    <div id="after-transaction-body">
                        <p>
                            ${this.playerStockQuantity + this.addPlayerStockQuantity} Shares = £${(this.playerStockQuantity + this.addPlayerStockQuantity) * this.gameStockValue}
                        </p>
                        <p>
                            ${this.gameStockSharesAvailable - this.addPlayerStockQuantity} Stocks left to buy
                        </p>
                        <p>
                            £${this.playerCash - this.addPlayerStockQuantity} Cash available
                        </p>
                    </div>
                </div>
            </div>
        `
    }

    connectedCallback() {

        this.querySelectorAll('.rem-stock').forEach(rem_stock_button => {
            rem_stock_button.addEventListener('click', () => {
                if (this.addPlayerStockQuantity >= 1000) {
                    this.addPlayerStockQuantity -= 1000
                }
            })
        })

        this.querySelectorAll('.add-stock').forEach(add_stock_button => {
            add_stock_button.addEventListener('click', () => {
                this.addPlayerStockQuantity += 1000
            })
        })

        /** @type {HTMLInputElement?} */
        let input_stock_count = this.querySelector('input.input-stock')

        if (input_stock_count == null) {
            throw new Error("input-stock element is required")
        }

        input_stock_count.value = this.addPlayerStockQuantity.toString()

        input_stock_count.addEventListener("change", (e) => {
            /** @type {HTMLInputElement?} */
            // @ts-ignore
            const target = e.target

            if (target == null || target.value == null) {
                throw new Error("input.input-stock event listener failed ")
            }

            this.addPlayerStockQuantity = parseInt(target.value)
        })

        /** @type {HTMLInputElement?} */
        let input_total_cost = this.querySelector('input.input-cash')

        if (input_total_cost == null) {
            throw new Error("input-cash element is required")
        }

        input_total_cost.setAttribute("step", String(1000 * this.gameStockValue))

        input_total_cost.addEventListener("change", (e) => {
            // alias for setting addPlayerStockQuantity according to the stockvalue

            /** @type {HTMLInputElement?} */
            // @ts-ignore
            const target = e.target

            if (target == null || target.value == null) {
                throw new Error("input.input-cash event listener failed ")
            }

            this.totalCost = parseInt(target.value)
        })

        const target_body_elem = this.querySelector('.transaction-body')

        if (target_body_elem == null) {
            throw new Error("transaction-body element is required")
        }

        target_body_elem.outerHTML = this.transaction_body

    }
}

customElements.define('stock-calculator', StockCalculator);