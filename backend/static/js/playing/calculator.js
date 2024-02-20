const html = String.raw

export class StockCalculator extends HTMLElement {

    static get observedAttributes() { 
        return [
            "stock-value",
            "stock-count",
            "playerStock-quantity", //={fmt.Sprintf("%d", player_stock.Quantity)}
            "playerStock-value",
            "gameStock-sharesAvailable", //={strconv.Itoa(player_stock.GameStock.SharesAvailable())}
            "player-cash" //={strconv.Itoa(player_stock.Player.Cash)}
        ]
    }

    set stockValue(value){
        this.setAttribute('stock-value', value)
        this.querySelector('.transaction-body').outerHTML = this.transaction_body
    }

    get stockValue(){
        return parseFloat(this.getAttribute('stock-value'))
    }

    set addPlayerStockQuantity(value){
        let rounded_stock_count = Math.floor(value / 1000) * 1000
        
        let input_stock_count = this.querySelector('.input-stock')
        input_stock_count.value = rounded_stock_count

        let input_total_cost = this.querySelector('.input-cash')
        input_total_cost.value = rounded_stock_count * this.stockValue
        
        this.setAttribute('stock-count', rounded_stock_count)

        this.querySelector('.transaction-body').outerHTML = this.transaction_body

        return rounded_stock_count
    }

    get addPlayerStockQuantity(){
        return parseInt(this.getAttribute('stock-count'))
    }

    set totalCost(value){

        let rounded_to_value_stock_count = Math.floor(value / this.stockValue);
        let new_stock_count = Math.floor(rounded_to_value_stock_count / 1000) * 1000
        this.addPlayerStockQuantity = new_stock_count

    }
    get totalCost(){
        return this.addPlayerStockQuantity * this.stockValue
    }

    get playerStockQuantity(){
        return parseInt(this.getAttribute('playerstock-quantity'))
    }

    get playerStockValue(){
        // float
        return parseFloat(this.getAttribute('playerstock-value'))
    }

    get gameStockSharesAvailable(){
        return parseInt(this.getAttribute('gamestock-sharesavailable'))
    }

    get playerCash(){
        return parseInt(this.getAttribute('player-cash'))
    }

    constructor() {
        super();
    }

    get transaction_body(){
        console.log("rerendering transaction body");
        return html`
            <div class="transaction-body" style="display: flex; flex-direction: row;">
                <div style="flex: 1;">
                    <h6>
                        Current
                    </h6>
                    
                    <p>
                        ${this.playerStockQuantity} Shares = £${this.playerStockQuantity * this.stockValue}
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
                            ${this.playerStockQuantity + this.addPlayerStockQuantity} Shares = £${(this.playerStockQuantity + this.addPlayerStockQuantity) * this.stockValue}
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

    attributeChangedCallback(name, oldValue, newValue) {}

    connectedCallback() {

        this.querySelectorAll('.rem-stock').forEach(rem_stock_button => rem_stock_button.addEventListener('click', (e) => {
            if(this.addPlayerStockQuantity >= 1000) { 
                this.addPlayerStockQuantity -= 1000
            }
            
        }))

        this.querySelectorAll('.add-stock').forEach(add_stock_button => add_stock_button.addEventListener('click', (e) => {
            this.addPlayerStockQuantity += 1000
        }))

        let input_stock_count = this.querySelector('.input-stock')
        input_stock_count.value = this.addPlayerStockQuantity
        
        input_stock_count.addEventListener("change", (e) => {
            this.addPlayerStockQuantity = e.target.value
        })

        let input_total_cost = this.querySelector('.input-cash')
        input_total_cost.setAttribute("step", 1000 * this.stockValue)
        
        input_total_cost.addEventListener("change", (e) => {
            // alias for setting addPlayerStockQuantity according to the stockvalue
            this.totalCost = e.target.value
        })

        this.querySelector('.transaction-body').outerHTML = this.transaction_body

    }
}

customElements.define('stock-calculator', StockCalculator);