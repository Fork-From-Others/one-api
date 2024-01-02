import React, {useEffect, useState} from 'react';
import {Button, Card, Grid, Header, Image, Modal, Pagination, Segment, Table} from 'semantic-ui-react';
import {API, showError} from "../../helpers";
import {PRODUCT_PER_PAGE} from "../../constants";

const Product = () => {
    const [products, setProducts] = useState([]);
    const [loading, setLoading] = useState(true);
    const [activePage, setActivePage] = useState(1);

    const [modalOpen, setModalOpen] = useState(false);
    const [qrcodeUrl, setQrcodeUrl] = useState(null);
    const [productId, setProductId] = useState(null);
    const [productKey, setProductKey] = useState(null);
    const [orderStatus, setOrderStatus] = useState('等待支付');

    const loadRedemptions = async (startIdx) => {
        let res = await API.get(`/api/redemption/pageQueryAndGroupBy?p=${startIdx}`);
        console.log('res: ', startIdx, res);
        const {success, message, data} = res.data;
        if (success) {
            if (startIdx === 0) {
                setProducts(data);
            } else {
                let newProducts = products;
                newProducts.push(...data);
                setProducts(newProducts);
            }
        } else {
            showError(message);
        }
        setLoading(false);
    }

    const onPaginationChange = (e, {activePage}) => {
        (async () => {
            if (activePage === Math.ceil(products.length / PRODUCT_PER_PAGE) + 1) {
                // In this case we have to load more data and then append them.
                await loadRedemptions(activePage - 1);
            }
            setActivePage(activePage);
        })();
    };

    useEffect(() => {
        loadRedemptions(0)
            .then()
            .catch((reason) => {
                showError(reason);
            });
    }, []);

    return (
        <>
            <Segment>
                <Header as='h3'>购买额度</Header>
                <Table basic compact size='small'>
                    <Grid columns={4} stackable>
                        {products.slice((activePage - 1) * PRODUCT_PER_PAGE, activePage * PRODUCT_PER_PAGE)
                            .map((item, index) => (
                                <Grid.Column key={index}>
                                    <Card fluid>
                                        <Card.Content>
                                            <Card.Header>{item.name}</Card.Header>
                                            <Card.Description>
                                                <p>额度: {item.quota}</p>
                                                <p>价格: {item.price}￥</p>
                                                <p></p>
                                            </Card.Description>
                                            <Button size={"mini"} positive fluid onClick={
                                                async () => {
                                                    let res = await API.post("/api/alipay/getPaymentQrcode", {
                                                        product_id: item.id,
                                                        subject: item.name,
                                                        money: item.price
                                                    });
                                                    const {success, message, data} = res.data;
                                                    const {product_id, out_trade_no, qr_code_url, qr_code_base64} = data;
                                                    setProductId(product_id);
                                                    if (success) {
                                                        console.log('data: ', data);
                                                        setQrcodeUrl(qr_code_base64);
                                                        setModalOpen(true);
                                                        // 开始定时查询订单状态
                                                        const intervalId = setInterval(async () => {
                                                            const res = await API.get(`/api/alipay/status?product_id=${product_id}&out_trade_no=${out_trade_no}`);
                                                            const {success, message, data} = res.data;
                                                            console.log('data: ', data);
                                                            if (data.status === '支付成功') {
                                                                setOrderStatus('支付成功');
                                                                setProductKey(data.key);
                                                                clearInterval(intervalId);
                                                            }
                                                        }, 5000);
                                                    } else {
                                                        showError(message);
                                                    }
                                                }
                                            }>下单</Button>
                                        </Card.Content>
                                    </Card>
                                </Grid.Column>
                            ))}
                    </Grid>
                    <Table.Footer>
                        <Pagination
                            activePage={activePage}
                            onPageChange={onPaginationChange}
                            size='mini'
                            siblingRange={1}
                            totalPages={
                                Math.ceil(products.length / PRODUCT_PER_PAGE) +
                                (products.length % PRODUCT_PER_PAGE === 0 ? 1 : 0)
                            }
                        />
                    </Table.Footer>
                </Table>
            </Segment>

            <Modal size={'mini'} open={modalOpen} onClose={() => setModalOpen(false)}>
                <Modal.Header>请扫描以下二维码进行支付</Modal.Header>
                <Modal.Content>
                    <Image size='medium' src={qrcodeUrl} wrapped alt="支付二维码"/>
                    <p>产品编号：{productId}</p>
                    <p>订单状态：{orderStatus}</p>
                    <p>兑换密钥：{productKey}</p>
                    <p style={{color: 'red'}}>当兑换密钥出现后请自行保存，不会再次显示！</p>
                </Modal.Content>
            </Modal>
        </>
    )
};

export default Product;