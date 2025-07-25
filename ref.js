async function createWidget() {
    const widget = new ListWidget();
    widget.backgroundColor = new Color('#1E1E1E');

    // 创建水平堆栈来分隔左右内容
    let mainStack = widget.addStack();
    mainStack.layoutHorizontally();

    // 左侧信息栈
    let leftStack = mainStack.addStack();
    leftStack.layoutVertically();
    leftStack.size = new Size(150, 150);

    // 设置请求头
    let headers = {
        Authorization: 'Bearer ',
        'User-Agent': 'Apifox/1.0.0 (https://apifox.com)',
        Accept: '*/*',
        Host: 'tapi.zeehoev.com',
        Connection: 'keep-alive',
        Cookie: 'acw_tc=0b32824217388280172008957ec68b4a84c95e1b5efd8b103d6c69b40480d9',
    };

    try {
        let url = 'https://tapi.zeehoev.com/v1.0/app/cfmotoserverapp/vehicle/widgets/358122500002456';
        let req = new Request(url);
        req.headers = headers;
        let response = await req.loadJSON();
        let data = response.data;

        // 添加车辆名称
        let nameText = leftStack.addText(data.vehicleName);
        nameText.font = Font.boldSystemFont(16);
        nameText.textColor = Color.white();
        leftStack.addSpacer(8);

        // 添加电量信息
        let batteryText = leftStack.addText(`电量: ${data.bmssoc}%`);
        batteryText.font = Font.systemFont(14);
        batteryText.textColor = Color.white();
        leftStack.addSpacer(4);

        // 添加续航里程
        let rangeText = leftStack.addText(`续航: ${data.hmiRidableMile}km`);
        rangeText.font = Font.systemFont(14);
        rangeText.textColor = Color.white();
        leftStack.addSpacer(4);

        // 添加位置信息
        if (data.location) {
            let locationText = leftStack.addText(`更新时间: ${data.location.locationTime}`);
            locationText.font = Font.systemFont(12);
            locationText.textColor = new Color('#888888');
        }

        // 右侧图片栈
        let rightStack = mainStack.addStack();
        rightStack.layoutVertically();
        rightStack.addSpacer(); // 将图片推到底部

        // 添加车辆图片
        let imageUrl = data.vehiclePicUrl;
        let imgReq = new Request(imageUrl);
        let img = await imgReq.loadImage();
        let imgWidget = rightStack.addImage(img);
        imgWidget.imageSize = new Size(150, 150); // 调整图片大小
        imgWidget.rightAlignImage(); // 右对齐图片
    } catch (err) {
        let errorText = widget.addText('获取数据失败');
        errorText.textColor = Color.red();
    }

    return widget;
}

// 运行小组件
let widget = await createWidget();
if (config.runsInWidget) {
    Script.setWidget(widget);
} else {
    widget.presentMedium();
}
Script.complete();
