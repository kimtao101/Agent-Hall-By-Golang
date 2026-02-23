package main

import (
	"fmt"
	"strings"
)

// XiaohongshuService 小红书文案生成服务
type XiaohongshuService struct {
	openai *OpenAIService
	logger Logger
}

// CopyRequest 文案生成请求
type CopyRequest struct {
	Scene  string                 `json:"scene" binding:"required"`
	Config map[string]interface{} `json:"config" binding:"required"`
}

// CopyResponse 文案生成响应
type CopyResponse struct {
	Copy string `json:"copy"`
}

// NewXiaohongshuService 创建新的小红书文案生成服务
func NewXiaohongshuService(openai *OpenAIService, logger Logger) *XiaohongshuService {
	if openai == nil {
		openai = GetDefaultOpenAIService()
	}
	if logger == nil {
		logger = NewDefaultLogger()
	}
	return &XiaohongshuService{
		openai: openai,
		logger: logger,
	}
}

// GenerateCopy 生成小红书文案
func (s *XiaohongshuService) GenerateCopy(req CopyRequest) (*CopyResponse, error) {
	s.logger.Info("Generating Xiaohongshu copy", map[string]interface{}{
		"scene":  req.Scene,
		"config": req.Config,
	})

	prompt := s.getPromptForScene(req.Scene, req.Config)
	s.logger.Info("Generated prompt", map[string]interface{}{
		"prompt": truncateString(prompt, 200) + "...",
	})

	messages := []Message{
		{
			Role:    "system",
			Content: "你是一位专业的小红书文案撰写专家，擅长撰写各种场景的优质文案。请根据用户提供的信息，生成符合小红书平台风格的文案，包含适当的表情符号、话题标签和互动引导语。",
		},
		{
			Role:    "user",
			Content: prompt,
		},
	}

	chatReq := ChatCompletionRequest{
		Messages:    messages,
		Model:       "deepseek-chat",
		Temperature: 0.8,
		Stream:      false,
	}

	resp, err := s.openai.CreateChatCompletion(chatReq)
	if err != nil {
		s.logger.Error("Error generating Xiaohongshu copy", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, err
	}

	var copy string
	if len(resp.Choices) > 0 {
		copy = resp.Choices[0].Message.Content
	}

	s.logger.Info("Generated copy successfully", map[string]interface{}{
		"copy": truncateString(copy, 200) + "...",
	})

	return &CopyResponse{Copy: copy}, nil
}

// getPromptForScene 根据场景生成提示词
func (s *XiaohongshuService) getPromptForScene(scene string, config map[string]interface{}) string {
	switch scene {
	case "beauty":
		return s.getBeautyPrompt(config)
	case "fashion":
		return s.getFashionPrompt(config)
	case "travel":
		return s.getTravelPrompt(config)
	case "food":
		return s.getFoodPrompt(config)
	case "home":
		return s.getHomePrompt(config)
	case "fitness":
		return s.getFitnessPrompt(config)
	case "parenting":
		return s.getParentingPrompt(config)
	case "tech":
		return s.getTechPrompt(config)
	default:
		return "请生成一篇小红书风格的文案"
	}
}

// getStringValue 从 map 中获取字符串值，如果不存在返回默认值
func getStringValue(config map[string]interface{}, key string, defaultValue string) string {
	if val, ok := config[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

// getBeautyPrompt 美妆护肤评测文案提示词
func (s *XiaohongshuService) getBeautyPrompt(config map[string]interface{}) string {
	productName := getStringValue(config, "productName", "")
	brand := getStringValue(config, "brand", "")
	price := getStringValue(config, "price", "未提供")
	skinType := getStringValue(config, "skinType", "未提供")
	texture := getStringValue(config, "texture", "未提供")
	keyIngredients := getStringValue(config, "keyIngredients", "未提供")
	usageFeel := getStringValue(config, "usageFeel", "")
	effect := getStringValue(config, "effect", "")
	recommendation := getStringValue(config, "recommendation", "")

	var prompt strings.Builder
	prompt.WriteString("请为以下美妆产品生成一篇小红书风格的评测文案：\n\n")
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", productName))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", brand))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", price))
	prompt.WriteString(fmt.Sprintf("适合肤质：%s\n", skinType))
	prompt.WriteString(fmt.Sprintf("质地：%s\n", texture))
	prompt.WriteString(fmt.Sprintf("核心成分：%s\n", keyIngredients))
	prompt.WriteString(fmt.Sprintf("使用感受：%s\n", usageFeel))
	prompt.WriteString(fmt.Sprintf("效果：%s\n", effect))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", recommendation))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 产品介绍（品牌、名称、价格等基本信息）\n")
	prompt.WriteString("3. 外观包装评价\n")
	prompt.WriteString("4. 质地和使用感受\n")
	prompt.WriteString("5. 效果评价\n")
	prompt.WriteString("6. 适合人群\n")
	prompt.WriteString("7. 总结推荐\n")
	prompt.WriteString("8. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格亲切自然，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语（如\"你们用过吗？\"、\"欢迎在评论区分享\"等）\n")
	prompt.WriteString("- 突出产品的核心卖点\n")
	prompt.WriteString("- 内容真实可信，避免过度夸大")

	return prompt.String()
}

// getFashionPrompt 穿搭搭配分享文案提示词
func (s *XiaohongshuService) getFashionPrompt(config map[string]interface{}) string {
	clothingType := getStringValue(config, "clothingType", "")
	style := getStringValue(config, "style", "")
	brand := getStringValue(config, "brand", "未提供")
	price := getStringValue(config, "price", "未提供")
	color := getStringValue(config, "color", "未提供")
	material := getStringValue(config, "material", "未提供")
	fit := getStringValue(config, "fit", "未提供")
	matchingTips := getStringValue(config, "matchingTips", "")
	scenario := getStringValue(config, "scenario", "")
	usageFeel := getStringValue(config, "usageFeel", "")

	var prompt strings.Builder
	prompt.WriteString("请为以下穿搭搭配生成一篇小红书风格的分享文案：\n\n")
	prompt.WriteString(fmt.Sprintf("服装类型：%s\n", clothingType))
	prompt.WriteString(fmt.Sprintf("风格：%s\n", style))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", brand))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", price))
	prompt.WriteString(fmt.Sprintf("颜色：%s\n", color))
	prompt.WriteString(fmt.Sprintf("材质：%s\n", material))
	prompt.WriteString(fmt.Sprintf("版型：%s\n", fit))
	prompt.WriteString(fmt.Sprintf("搭配建议：%s\n", matchingTips))
	prompt.WriteString(fmt.Sprintf("适合场景：%s\n", scenario))
	prompt.WriteString(fmt.Sprintf("穿着感受：%s\n\n", usageFeel))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 整体搭配介绍\n")
	prompt.WriteString("3. 单品推荐（材质、版型等）\n")
	prompt.WriteString("4. 搭配技巧和思路\n")
	prompt.WriteString("5. 适合场景和人群\n")
	prompt.WriteString("6. 穿着感受\n")
	prompt.WriteString("7. 购买建议\n")
	prompt.WriteString("8. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格时尚潮流，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语\n")
	prompt.WriteString("- 提供具体的搭配建议\n")
	prompt.WriteString("- 突出服装的风格特点")

	return prompt.String()
}

// getTravelPrompt 旅行打卡攻略文案提示词
func (s *XiaohongshuService) getTravelPrompt(config map[string]interface{}) string {
	destination := getStringValue(config, "destination", "")
	duration := getStringValue(config, "duration", "")
	bestTime := getStringValue(config, "bestTime", "未提供")
	budget := getStringValue(config, "budget", "未提供")
	attractions := getStringValue(config, "attractions", "")
	food := getStringValue(config, "food", "")
	accommodation := getStringValue(config, "accommodation", "未提供")
	transportation := getStringValue(config, "transportation", "未提供")
	tips := getStringValue(config, "tips", "")
	experience := getStringValue(config, "experience", "")

	var prompt strings.Builder
	prompt.WriteString("请为以下旅行目的地生成一篇小红书风格的打卡攻略文案：\n\n")
	prompt.WriteString(fmt.Sprintf("目的地：%s\n", destination))
	prompt.WriteString(fmt.Sprintf("行程天数：%s天\n", duration))
	prompt.WriteString(fmt.Sprintf("最佳时间：%s\n", bestTime))
	prompt.WriteString(fmt.Sprintf("预算：%s\n", budget))
	prompt.WriteString(fmt.Sprintf("主要景点：%s\n", attractions))
	prompt.WriteString(fmt.Sprintf("特色美食：%s\n", food))
	prompt.WriteString(fmt.Sprintf("住宿推荐：%s\n", accommodation))
	prompt.WriteString(fmt.Sprintf("交通方式：%s\n", transportation))
	prompt.WriteString(fmt.Sprintf("旅行贴士：%s\n", tips))
	prompt.WriteString(fmt.Sprintf("个人体验：%s\n\n", experience))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 旅行概览（目的地、天数、预算等）\n")
	prompt.WriteString("3. 行程安排推荐\n")
	prompt.WriteString("4. 景点打卡攻略\n")
	prompt.WriteString("5. 美食推荐\n")
	prompt.WriteString("6. 住宿和交通建议\n")
	prompt.WriteString("7. 实用贴士\n")
	prompt.WriteString("8. 个人感受和总结\n")
	prompt.WriteString("9. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格轻松愉快，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语\n")
	prompt.WriteString("- 提供详细的实用信息\n")
	prompt.WriteString("- 突出目的地的特色和亮点")

	return prompt.String()
}

// getFoodPrompt 美食探店体验文案提示词
func (s *XiaohongshuService) getFoodPrompt(config map[string]interface{}) string {
	restaurantName := getStringValue(config, "restaurantName", "")
	location := getStringValue(config, "location", "")
	cuisineType := getStringValue(config, "cuisineType", "")
	priceRange := getStringValue(config, "priceRange", "未提供")
	environment := getStringValue(config, "environment", "未提供")
	service := getStringValue(config, "service", "未提供")
	signatureDishes := getStringValue(config, "signatureDishes", "")
	taste := getStringValue(config, "taste", "")
	recommendation := getStringValue(config, "recommendation", "")

	var prompt strings.Builder
	prompt.WriteString("请为以下餐厅生成一篇小红书风格的探店体验文案：\n\n")
	prompt.WriteString(fmt.Sprintf("餐厅名称：%s\n", restaurantName))
	prompt.WriteString(fmt.Sprintf("位置：%s\n", location))
	prompt.WriteString(fmt.Sprintf("菜系：%s\n", cuisineType))
	prompt.WriteString(fmt.Sprintf("价格区间：%s\n", priceRange))
	prompt.WriteString(fmt.Sprintf("环境：%s\n", environment))
	prompt.WriteString(fmt.Sprintf("服务：%s\n", service))
	prompt.WriteString(fmt.Sprintf("招牌菜：%s\n", signatureDishes))
	prompt.WriteString(fmt.Sprintf("口味：%s\n", taste))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", recommendation))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 餐厅基本信息（位置、环境等）\n")
	prompt.WriteString("3. 招牌菜推荐和评价\n")
	prompt.WriteString("4. 口味和服务评价\n")
	prompt.WriteString("5. 价格和性价比\n")
	prompt.WriteString("6. 适合人群和场景\n")
	prompt.WriteString("7. 总结推荐\n")
	prompt.WriteString("8. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格生动诱人，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语\n")
	prompt.WriteString("- 提供详细的菜品评价\n")
	prompt.WriteString("- 突出餐厅的特色和亮点")

	return prompt.String()
}

// getHomePrompt 家居好物推荐文案提示词
func (s *XiaohongshuService) getHomePrompt(config map[string]interface{}) string {
	productName := getStringValue(config, "productName", "")
	category := getStringValue(config, "category", "")
	brand := getStringValue(config, "brand", "未提供")
	price := getStringValue(config, "price", "未提供")
	material := getStringValue(config, "material", "未提供")
	size := getStringValue(config, "size", "未提供")
	usageScenario := getStringValue(config, "usageScenario", "")
	functionality := getStringValue(config, "functionality", "")
	usageFeel := getStringValue(config, "usageFeel", "")
	spaceSaving := getStringValue(config, "spaceSaving", "未提供")
	recommendation := getStringValue(config, "recommendation", "")

	var prompt strings.Builder
	prompt.WriteString("请为以下家居好物生成一篇小红书风格的推荐文案：\n\n")
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", productName))
	prompt.WriteString(fmt.Sprintf("类别：%s\n", category))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", brand))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", price))
	prompt.WriteString(fmt.Sprintf("材质：%s\n", material))
	prompt.WriteString(fmt.Sprintf("尺寸：%s\n", size))
	prompt.WriteString(fmt.Sprintf("使用场景：%s\n", usageScenario))
	prompt.WriteString(fmt.Sprintf("功能：%s\n", functionality))
	prompt.WriteString(fmt.Sprintf("使用感受：%s\n", usageFeel))
	prompt.WriteString(fmt.Sprintf("节省空间：%s\n", spaceSaving))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", recommendation))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 产品介绍（名称、类别、价格等）\n")
	prompt.WriteString("3. 外观设计评价\n")
	prompt.WriteString("4. 功能和使用方法\n")
	prompt.WriteString("5. 使用感受和效果\n")
	prompt.WriteString("6. 适用场景\n")
	prompt.WriteString("7. 性价比评价\n")
	prompt.WriteString("8. 总结推荐\n")
	prompt.WriteString("9. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格实用亲切，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语\n")
	prompt.WriteString("- 提供详细的使用体验\n")
	prompt.WriteString("- 突出产品的实用价值")

	return prompt.String()
}

// getFitnessPrompt 健身运动记录文案提示词
func (s *XiaohongshuService) getFitnessPrompt(config map[string]interface{}) string {
	workoutType := getStringValue(config, "workoutType", "")
	duration := getStringValue(config, "duration", "未提供")
	frequency := getStringValue(config, "frequency", "未提供")
	equipment := getStringValue(config, "equipment", "未提供")
	difficulty := getStringValue(config, "difficulty", "未提供")
	benefits := getStringValue(config, "benefits", "")
	experience := getStringValue(config, "experience", "")
	tips := getStringValue(config, "tips", "")
	results := getStringValue(config, "results", "未提供")

	var prompt strings.Builder
	prompt.WriteString("请为以下健身运动生成一篇小红书风格的记录文案：\n\n")
	prompt.WriteString(fmt.Sprintf("运动类型：%s\n", workoutType))
	prompt.WriteString(fmt.Sprintf("运动时长：%s分钟\n", duration))
	prompt.WriteString(fmt.Sprintf("运动频率：%s\n", frequency))
	prompt.WriteString(fmt.Sprintf("所需装备：%s\n", equipment))
	prompt.WriteString(fmt.Sprintf("难度：%s\n", difficulty))
	prompt.WriteString(fmt.Sprintf("运动好处：%s\n", benefits))
	prompt.WriteString(fmt.Sprintf("个人体验：%s\n", experience))
	prompt.WriteString(fmt.Sprintf("注意事项：%s\n", tips))
	prompt.WriteString(fmt.Sprintf("运动效果：%s\n\n", results))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 运动介绍（类型、时长、频率等）\n")
	prompt.WriteString("3. 运动过程和感受\n")
	prompt.WriteString("4. 运动好处和效果\n")
	prompt.WriteString("5. 适合人群\n")
	prompt.WriteString("6. 注意事项和建议\n")
	prompt.WriteString("7. 个人心得和激励\n")
	prompt.WriteString("8. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格积极向上，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语\n")
	prompt.WriteString("- 提供详细的运动体验\n")
	prompt.WriteString("- 传递正能量和激励信息")

	return prompt.String()
}

// getParentingPrompt 母婴育儿心得文案提示词
func (s *XiaohongshuService) getParentingPrompt(config map[string]interface{}) string {
	babyAge := getStringValue(config, "babyAge", "")
	topic := getStringValue(config, "topic", "")
	productName := getStringValue(config, "productName", "未提供")
	brand := getStringValue(config, "brand", "未提供")
	price := getStringValue(config, "price", "未提供")
	problem := getStringValue(config, "problem", "")
	solution := getStringValue(config, "solution", "")
	experience := getStringValue(config, "experience", "")
	tips := getStringValue(config, "tips", "")

	var prompt strings.Builder
	prompt.WriteString("请为以下母婴育儿主题生成一篇小红书风格的心得分享文案：\n\n")
	prompt.WriteString(fmt.Sprintf("宝宝年龄：%s\n", babyAge))
	prompt.WriteString(fmt.Sprintf("主题：%s\n", topic))
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", productName))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", brand))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", price))
	prompt.WriteString(fmt.Sprintf("问题描述：%s\n", problem))
	prompt.WriteString(fmt.Sprintf("解决方案：%s\n", solution))
	prompt.WriteString(fmt.Sprintf("经验分享：%s\n", experience))
	prompt.WriteString(fmt.Sprintf("小贴士：%s\n\n", tips))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 问题背景介绍\n")
	prompt.WriteString("3. 解决方案分享\n")
	prompt.WriteString("4. 经验总结\n")
	prompt.WriteString("5. 实用小贴士\n")
	prompt.WriteString("6. 产品推荐（如果有）\n")
	prompt.WriteString("7. 互动和鼓励\n")
	prompt.WriteString("8. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格温暖亲切，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语\n")
	prompt.WriteString("- 提供实用的育儿经验\n")
	prompt.WriteString("- 传递正能量和支持信息")

	return prompt.String()
}

// getTechPrompt 数码产品测评文案提示词
func (s *XiaohongshuService) getTechPrompt(config map[string]interface{}) string {
	productName := getStringValue(config, "productName", "")
	brand := getStringValue(config, "brand", "")
	price := getStringValue(config, "price", "未提供")
	releaseDate := getStringValue(config, "releaseDate", "未提供")
	specs := getStringValue(config, "specs", "")
	design := getStringValue(config, "design", "未提供")
	performance := getStringValue(config, "performance", "")
	battery := getStringValue(config, "battery", "未提供")
	camera := getStringValue(config, "camera", "未提供")
	userExperience := getStringValue(config, "userExperience", "")
	pros := getStringValue(config, "pros", "")
	cons := getStringValue(config, "cons", "未提供")
	recommendation := getStringValue(config, "recommendation", "")

	var prompt strings.Builder
	prompt.WriteString("请为以下数码产品生成一篇小红书风格的测评文案：\n\n")
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", productName))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", brand))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", price))
	prompt.WriteString(fmt.Sprintf("上市时间：%s\n", releaseDate))
	prompt.WriteString(fmt.Sprintf("主要配置：%s\n", specs))
	prompt.WriteString(fmt.Sprintf("外观设计：%s\n", design))
	prompt.WriteString(fmt.Sprintf("性能表现：%s\n", performance))
	prompt.WriteString(fmt.Sprintf("续航：%s\n", battery))
	prompt.WriteString(fmt.Sprintf("相机：%s\n", camera))
	prompt.WriteString(fmt.Sprintf("用户体验：%s\n", userExperience))
	prompt.WriteString(fmt.Sprintf("优点：%s\n", pros))
	prompt.WriteString(fmt.Sprintf("缺点：%s\n", cons))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", recommendation))
	prompt.WriteString("请按照以下结构撰写：\n")
	prompt.WriteString("1. 吸引人的标题（包含emoji）\n")
	prompt.WriteString("2. 产品开箱介绍\n")
	prompt.WriteString("3. 外观设计评价\n")
	prompt.WriteString("4. 性能和功能测试\n")
	prompt.WriteString("5. 用户体验\n")
	prompt.WriteString("6. 优缺点分析\n")
	prompt.WriteString("7. 适合人群\n")
	prompt.WriteString("8. 性价比评价\n")
	prompt.WriteString("9. 总结推荐\n")
	prompt.WriteString("10. 相关话题标签（至少5个）\n\n")
	prompt.WriteString("要求：\n")
	prompt.WriteString("- 语言风格专业易懂，符合小红书用户的阅读习惯\n")
	prompt.WriteString("- 适当使用emoji表情符号\n")
	prompt.WriteString("- 包含互动引导语\n")
	prompt.WriteString("- 提供详细的产品信息\n")
	prompt.WriteString("- 评价客观公正，避免过度 bias")

	return prompt.String()
}
