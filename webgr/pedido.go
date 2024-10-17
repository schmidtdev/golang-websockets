package webgr

type WebgrPedido struct {
	Cnpedido                   int    `gorm:"primaryKey;autoIncrement" json:"cnpedido"`
	Nrpedido                   int    `json:"nrpedido"`
	Cdpedido                   int    `json:"cdpedido"`
	Sufixo                     string `gorm:"size:3" json:"sufixo"`
	Nrpedidoerp                int    `json:"nrpedidoerp"`
	Nmcontato                  string `gorm:"size:50" json:"nmcontato"`
	Xped                       string `gorm:"size:60" json:"xped"`
	Cnrepresentante            int    `gorm:"not null" json:"cnrepresentante"`
	Cncliente                  int    `gorm:"not null" json:"cncliente"`
	Cnclientebenefeciario      int    `json:"cnclientebenefeciario"`
	Cntpoperacaosaida          int    `gorm:"not null" json:"cntpoperacaosaida"`
	Cntipopedido               int    `json:"cntipopedido"`
	Cnusuario                  int    `gorm:"not null" json:"cnusuario"`
	Cnformapgto                int    `json:"cnformapgto"`
	Cncondicaopgto             int    `json:"cncondicaopgto"`
	Cnempresa                  int    `gorm:"not null" json:"cnempresa"`
	Cnvisita                   int    `json:"cnvisita"`
	Cntransportador            string `json:"cntransportador"`
	Obsinterna                 string `gorm:"size:960" json:"obsinterna"`
	Pedaviso                   string `gorm:"size:500" json:"pedaviso"`
	Cdaut                      string `gorm:"size:250" json:"cdaut"`
	Ftaut                      string `gorm:"type:text" json:"ftaut"`
	Cntblvenda                 int    `json:"cntblvenda"`
	Prtmudartblvenda           string `gorm:"size:2;default:'S'" json:"prtmudartblvenda"`
	Retira                     string `gorm:"size:2;default:'N'" json:"retira"`
	Vrbaseicms                 string `json:"vrbaseicms"`
	Porcentagemicms            string `json:"porcentagemicms"`
	Vricms                     string `json:"vricms"`
	Vrbasest                   string `json:"vrbasest"`
	Vrst                       string `json:"vrst"`
	Obspedido                  string `gorm:"size:960" json:"obspedido"`
	Cnempresaorigem            int    `json:"cnempresaorigem"`
	Isfixotpoperacaosaida      string `gorm:"size:2;default:'N'" json:"isfixotpoperacaosaida"`
	Isfixoformapgto            string `gorm:"size:2;default:'N'" json:"isfixoformapgto"`
	Isfixocondicaopgto         string `gorm:"size:2;default:'N'" json:"isfixocondicaopgto"`
	Isfixotblvenda             string `gorm:"size:2;default:'N'" json:"isfixotblvenda"`
	Status                     string `gorm:"size:2;default:'10';not null" json:"status"`
	Nrnfe                      string `gorm:"size:9" json:"nrnfe"`
	Nrcfe                      string `gorm:"size:20" json:"nrcfe"`
	Erpstatus                  string `gorm:"size:3" json:"erpstatus"`
	Motivocancelamento         string `gorm:"size:250" json:"motivocancelamento"`
	Erpstatusdesc              string `gorm:"size:50" json:"erpstatusdesc"`
	Latitude                   string `json:"latitude"`
	Longitude                  string `json:"longitude"`
	Ispedidooficial            string `gorm:"size:2;default:'N'" json:"ispedidooficial"`
	Percentualcomissao         string `gorm:"default:'0.00'" json:"percentualcomissao"`
	Vrcomissao                 string `gorm:"default:'0.0'" json:"vrcomissao"`
	Vrbruto                    string `gorm:"default:'0.0'" json:"vrbruto"`
	Vrdesconto                 string `gorm:"default:'0.0'" json:"vrdesconto"`
	Vrdescontototal            string `gorm:"default:'0.00'" json:"vrdescontototal"`
	Vrdescontofinal            string `gorm:"default:'0.00'" json:"vrdescontofinal"`
	Percentualdescontofinal    string `gorm:"default:0.00" json:"percentualdescontofinal"`
	Tipodesconto               string `gorm:"size:2;default:'P'" json:"tipodesconto"`
	Vrdescontocondpgto         string `gorm:"default:'0.00'" json:"vrdescontocondpgto"`
	Percentualdescontocondpgto string `gorm:"default:0.00" json:"percentualdescontocondpgto"`
	Tipodescontocondpgto       string `gorm:"size:2;default:'P'" json:"tipodescontocondpgto"`
	Vrliquidototal             string `gorm:"default:'0.00'" json:"vrliquidototal"`
	Vrtotal                    string `gorm:"default:'0.00'" json:"vrtotal"`
}
