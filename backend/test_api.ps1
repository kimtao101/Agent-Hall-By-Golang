# 小红书智能体后端 API 测试脚本

$baseUrl = "http://localhost:8016"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "小红书智能体后端 API 测试" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan

# 1. 健康检查
Write-Host "`n[1] 测试健康检查接口..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/health" -Method GET
    Write-Host "✓ 健康检查通过: $($response.status)" -ForegroundColor Green
    Write-Host "  时间戳: $($response.timestamp)"
} catch {
    Write-Host "✗ 健康检查失败: $_" -ForegroundColor Red
}

# 2. 测试小红书文案生成 - 美妆场景
Write-Host "`n[2] 测试小红书文案生成接口 (美妆场景)..." -ForegroundColor Yellow
try {
    $body = @{
        scene = "beauty"
        config = @{
            productName = "小棕瓶精华"
            brand = "雅诗兰黛"
            price = "999"
            skinType = "混合性"
            texture = "清爽"
            keyIngredients = "二裂酵母"
            usageFeel = "吸收很快，不油腻，上脸很舒服"
            effect = "肤色提亮明显，毛孔变小"
            recommendation = "适合熬夜肌，值得回购"
        }
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "$baseUrl/xiaohongshu/copy" -Method POST -ContentType "application/json" -Body $body
    Write-Host "✓ 文案生成成功" -ForegroundColor Green
    Write-Host "  生成的文案:" -ForegroundColor Gray
    Write-Host "  $($response.copy.Substring(0, [Math]::Min(200, $response.copy.Length)))..." -ForegroundColor Gray
} catch {
    Write-Host "✗ 文案生成失败: $_" -ForegroundColor Red
    if ($_.Exception.Response) {
        $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
        $reader.BaseStream.Position = 0
        $reader.DiscardBufferedData()
        $errorBody = $reader.ReadToEnd()
        Write-Host "  错误详情: $errorBody" -ForegroundColor Red
    }
}

# 3. 测试小红书文案生成 - 美食场景
Write-Host "`n[3] 测试小红书文案生成接口 (美食场景)..." -ForegroundColor Yellow
try {
    $body = @{
        scene = "food"
        config = @{
            restaurantName = "蜀大侠火锅"
            location = "春熙路"
            cuisineType = "川味火锅"
            priceRange = "人均120"
            signatureDishes = "毛肚、鸭肠、嫩牛肉"
            taste = "麻辣鲜香，锅底浓郁"
            recommendation = "服务热情，食材新鲜"
        }
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "$baseUrl/xiaohongshu/copy" -Method POST -ContentType "application/json" -Body $body
    Write-Host "✓ 文案生成成功" -ForegroundColor Green
    Write-Host "  生成的文案:" -ForegroundColor Gray
    Write-Host "  $($response.copy.Substring(0, [Math]::Min(200, $response.copy.Length)))..." -ForegroundColor Gray
} catch {
    Write-Host "✗ 文案生成失败: $_" -ForegroundColor Red
}

# 4. 测试聊天接口（流式）
Write-Host "`n[4] 测试聊天接口..." -ForegroundColor Yellow
try {
    $body = @{
        message = "你好，请介绍一下自己"
    } | ConvertTo-Json
    
    $response = Invoke-RestMethod -Uri "$baseUrl/chat" -Method POST -ContentType "application/json" -Body $body
    Write-Host "✓ 聊天接口调用成功" -ForegroundColor Green
    Write-Host "  响应: $($response.Substring(0, [Math]::Min(100, $response.Length)))..." -ForegroundColor Gray
} catch {
    Write-Host "✗ 聊天接口调用失败: $_" -ForegroundColor Red
}

# 5. 获取聊天历史
Write-Host "`n[5] 测试获取聊天历史接口..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$baseUrl/history" -Method GET
    Write-Host "✓ 获取历史成功，共 $($response.Count) 条消息" -ForegroundColor Green
} catch {
    Write-Host "✗ 获取历史失败: $_" -ForegroundColor Red
}

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "测试完成" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
