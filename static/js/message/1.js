	//从缓存中获取数据并渲染
	let msgBoxList = JSON.parse(window.localStorage.getItem('msgBoxList')) || [];
	innerHTMl(msgBoxList)


	//点击发表按扭，发表内容
	$("span.submit").click(function () {
		let txt = document.getElementById("messag").value;
		if (!txt) {
			$('.message').focus(); //自动获取焦点
			return;
		}
		let obj = {
			msg: txt
		}
		msgBoxList.unshift(obj) //添加到数组里
		window.localStorage.setItem('msgBoxList', JSON.stringify(msgBoxList)) //将数据保存到缓存
		innerHTMl([obj]) //渲染当前输入框内容
		$('.message').val('') // 清空输入框

	});


	//回复按钮
	$("body").on('click', '.huifu', function () {

		$(".huifu").click(function(){
                $(".reply_textarea").remove();
                $(this).parent().append("<div id='q1q' class='reply_textarea'><textarea class='test1'></textarea><div class='But'><br><span class='test'>发表</span></div></div>");
            });

	})
	
	//获取当前时间
	function getNowFormatDate() {
        var date = new Date();
        var seperator1 = "-";
        var year = date.getFullYear();
        var month = date.getMonth() + 1;
        var strDate = date.getDate();
        if (month >= 1 && month <= 9) {
            month = "0" + month;
        }
        if (strDate >= 0 && strDate <= 9) {
            strDate = "0" + strDate;
        }
        var currentdate = year + seperator1 + month + seperator1 + strDate;
        return currentdate;
    }

    //临时用户
	function name() {
        var namea = '马太阳';
        return namea;
    }

	//渲染html
	function innerHTMl(List) {
		List = List || []
		List.forEach(item => {
			let str =
				`<div class='msgBox'>
					<div class="headUrl">
						<div>
							<span class="time">用户:`+name()+`&nbsp;&nbsp;&nbsp;&nbsp;回复时间:`+getNowFormatDate()+`</span>
						</div>
						
					</div>
					<div class='msgTxt'>
						${item.msg}
					</div>
					<span class='huifu' href='javascript:;'>回复</span><br><br>
				</div>`
			$(".msgCon").prepend(str);
		})
	}