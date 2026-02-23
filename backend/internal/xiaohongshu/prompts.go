package xiaohongshu

import (
	"fmt"
	"strings"

	"agent-backend/pkg/utils"
)

// PromptBuilder 提示词构建器
type PromptBuilder struct {
	config map[string]interface{}
}

// NewPromptBuilder 创建新的提示词构建器
func NewPromptBuilder(config map[string]interface{}) *PromptBuilder {
	return &PromptBuilder{config: config}
}

// Build 根据场景构建提示词
func (pb *PromptBuilder) Build(scene string) string {
	switch scene {
	case "beauty":
		return pb.buildBeautyPrompt()
	case "fashion":
		return pb.buildFashionPrompt()
	case "travel":
		return pb.buildTravelPrompt()
	case "food":
		return pb.buildFoodPrompt()
	case "home":
		return pb.buildHomePrompt()
	case "fitness":
		return pb.buildFitnessPrompt()
	case "parenting":
		return pb.buildParentingPrompt()
	case "tech":
		return pb.buildTechPrompt()
	default:
		return "请生成一篇小红书风格的文案"
	}
}

// buildBeautyPrompt 美妆护肤评测文案提示词
func (pb *PromptBuilder) buildBeautyPrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下美妆产品生成一篇小红书风格的评测文案：\n\n")
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", utils.GetStringValue(c, "productName", "")))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", utils.GetStringValue(c, "brand", "")))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", utils.GetStringValue(c, "price", "未提供")))
	prompt.WriteString(fmt.Sprintf("适合肤质：%s\n", utils.GetStringValue(c, "skinType", "未提供")))
	prompt.WriteString(fmt.Sprintf("质地：%s\n", utils.GetStringValue(c, "texture", "未提供")))
	prompt.WriteString(fmt.Sprintf("核心成分：%s\n", utils.GetStringValue(c, "keyIngredients", "未提供")))
	prompt.WriteString(fmt.Sprintf("使用感受：%s\n", utils.GetStringValue(c, "usageFeel", "")))
	prompt.WriteString(fmt.Sprintf("效果：%s\n", utils.GetStringValue(c, "effect", "")))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", utils.GetStringValue(c, "recommendation", "")))
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

// buildFashionPrompt 穿搭搭配分享文案提示词
func (pb *PromptBuilder) buildFashionPrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下穿搭搭配生成一篇小红书风格的分享文案：\n\n")
	prompt.WriteString(fmt.Sprintf("服装类型：%s\n", utils.GetStringValue(c, "clothingType", "")))
	prompt.WriteString(fmt.Sprintf("风格：%s\n", utils.GetStringValue(c, "style", "")))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", utils.GetStringValue(c, "brand", "未提供")))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", utils.GetStringValue(c, "price", "未提供")))
	prompt.WriteString(fmt.Sprintf("颜色：%s\n", utils.GetStringValue(c, "color", "未提供")))
	prompt.WriteString(fmt.Sprintf("材质：%s\n", utils.GetStringValue(c, "material", "未提供")))
	prompt.WriteString(fmt.Sprintf("版型：%s\n", utils.GetStringValue(c, "fit", "未提供")))
	prompt.WriteString(fmt.Sprintf("搭配建议：%s\n", utils.GetStringValue(c, "matchingTips", "")))
	prompt.WriteString(fmt.Sprintf("适合场景：%s\n", utils.GetStringValue(c, "scenario", "")))
	prompt.WriteString(fmt.Sprintf("穿着感受：%s\n\n", utils.GetStringValue(c, "usageFeel", "")))
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

// buildTravelPrompt 旅行打卡攻略文案提示词
func (pb *PromptBuilder) buildTravelPrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下旅行目的地生成一篇小红书风格的打卡攻略文案：\n\n")
	prompt.WriteString(fmt.Sprintf("目的地：%s\n", utils.GetStringValue(c, "destination", "")))
	prompt.WriteString(fmt.Sprintf("行程天数：%s天\n", utils.GetStringValue(c, "duration", "")))
	prompt.WriteString(fmt.Sprintf("最佳时间：%s\n", utils.GetStringValue(c, "bestTime", "未提供")))
	prompt.WriteString(fmt.Sprintf("预算：%s\n", utils.GetStringValue(c, "budget", "未提供")))
	prompt.WriteString(fmt.Sprintf("主要景点：%s\n", utils.GetStringValue(c, "attractions", "")))
	prompt.WriteString(fmt.Sprintf("特色美食：%s\n", utils.GetStringValue(c, "food", "")))
	prompt.WriteString(fmt.Sprintf("住宿推荐：%s\n", utils.GetStringValue(c, "accommodation", "未提供")))
	prompt.WriteString(fmt.Sprintf("交通方式：%s\n", utils.GetStringValue(c, "transportation", "未提供")))
	prompt.WriteString(fmt.Sprintf("旅行贴士：%s\n", utils.GetStringValue(c, "tips", "")))
	prompt.WriteString(fmt.Sprintf("个人体验：%s\n\n", utils.GetStringValue(c, "experience", "")))
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

// buildFoodPrompt 美食探店体验文案提示词
func (pb *PromptBuilder) buildFoodPrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下餐厅生成一篇小红书风格的探店体验文案：\n\n")
	prompt.WriteString(fmt.Sprintf("餐厅名称：%s\n", utils.GetStringValue(c, "restaurantName", "")))
	prompt.WriteString(fmt.Sprintf("位置：%s\n", utils.GetStringValue(c, "location", "")))
	prompt.WriteString(fmt.Sprintf("菜系：%s\n", utils.GetStringValue(c, "cuisineType", "")))
	prompt.WriteString(fmt.Sprintf("价格区间：%s\n", utils.GetStringValue(c, "priceRange", "未提供")))
	prompt.WriteString(fmt.Sprintf("环境：%s\n", utils.GetStringValue(c, "environment", "未提供")))
	prompt.WriteString(fmt.Sprintf("服务：%s\n", utils.GetStringValue(c, "service", "未提供")))
	prompt.WriteString(fmt.Sprintf("招牌菜：%s\n", utils.GetStringValue(c, "signatureDishes", "")))
	prompt.WriteString(fmt.Sprintf("口味：%s\n", utils.GetStringValue(c, "taste", "")))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", utils.GetStringValue(c, "recommendation", "")))
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

// buildHomePrompt 家居好物推荐文案提示词
func (pb *PromptBuilder) buildHomePrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下家居好物生成一篇小红书风格的推荐文案：\n\n")
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", utils.GetStringValue(c, "productName", "")))
	prompt.WriteString(fmt.Sprintf("类别：%s\n", utils.GetStringValue(c, "category", "")))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", utils.GetStringValue(c, "brand", "未提供")))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", utils.GetStringValue(c, "price", "未提供")))
	prompt.WriteString(fmt.Sprintf("材质：%s\n", utils.GetStringValue(c, "material", "未提供")))
	prompt.WriteString(fmt.Sprintf("尺寸：%s\n", utils.GetStringValue(c, "size", "未提供")))
	prompt.WriteString(fmt.Sprintf("使用场景：%s\n", utils.GetStringValue(c, "usageScenario", "")))
	prompt.WriteString(fmt.Sprintf("功能：%s\n", utils.GetStringValue(c, "functionality", "")))
	prompt.WriteString(fmt.Sprintf("使用感受：%s\n", utils.GetStringValue(c, "usageFeel", "")))
	prompt.WriteString(fmt.Sprintf("节省空间：%s\n", utils.GetStringValue(c, "spaceSaving", "未提供")))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", utils.GetStringValue(c, "recommendation", "")))
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

// buildFitnessPrompt 健身运动记录文案提示词
func (pb *PromptBuilder) buildFitnessPrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下健身运动生成一篇小红书风格的记录文案：\n\n")
	prompt.WriteString(fmt.Sprintf("运动类型：%s\n", utils.GetStringValue(c, "workoutType", "")))
	prompt.WriteString(fmt.Sprintf("运动时长：%s分钟\n", utils.GetStringValue(c, "duration", "未提供")))
	prompt.WriteString(fmt.Sprintf("运动频率：%s\n", utils.GetStringValue(c, "frequency", "未提供")))
	prompt.WriteString(fmt.Sprintf("所需装备：%s\n", utils.GetStringValue(c, "equipment", "未提供")))
	prompt.WriteString(fmt.Sprintf("难度：%s\n", utils.GetStringValue(c, "difficulty", "未提供")))
	prompt.WriteString(fmt.Sprintf("运动好处：%s\n", utils.GetStringValue(c, "benefits", "")))
	prompt.WriteString(fmt.Sprintf("个人体验：%s\n", utils.GetStringValue(c, "experience", "")))
	prompt.WriteString(fmt.Sprintf("注意事项：%s\n", utils.GetStringValue(c, "tips", "")))
	prompt.WriteString(fmt.Sprintf("运动效果：%s\n\n", utils.GetStringValue(c, "results", "未提供")))
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

// buildParentingPrompt 母婴育儿心得文案提示词
func (pb *PromptBuilder) buildParentingPrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下母婴育儿主题生成一篇小红书风格的心得分享文案：\n\n")
	prompt.WriteString(fmt.Sprintf("宝宝年龄：%s\n", utils.GetStringValue(c, "babyAge", "")))
	prompt.WriteString(fmt.Sprintf("主题：%s\n", utils.GetStringValue(c, "topic", "")))
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", utils.GetStringValue(c, "productName", "未提供")))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", utils.GetStringValue(c, "brand", "未提供")))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", utils.GetStringValue(c, "price", "未提供")))
	prompt.WriteString(fmt.Sprintf("问题描述：%s\n", utils.GetStringValue(c, "problem", "")))
	prompt.WriteString(fmt.Sprintf("解决方案：%s\n", utils.GetStringValue(c, "solution", "")))
	prompt.WriteString(fmt.Sprintf("经验分享：%s\n", utils.GetStringValue(c, "experience", "")))
	prompt.WriteString(fmt.Sprintf("小贴士：%s\n\n", utils.GetStringValue(c, "tips", "")))
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

// buildTechPrompt 数码产品测评文案提示词
func (pb *PromptBuilder) buildTechPrompt() string {
	c := pb.config
	var prompt strings.Builder
	prompt.WriteString("请为以下数码产品生成一篇小红书风格的测评文案：\n\n")
	prompt.WriteString(fmt.Sprintf("产品名称：%s\n", utils.GetStringValue(c, "productName", "")))
	prompt.WriteString(fmt.Sprintf("品牌：%s\n", utils.GetStringValue(c, "brand", "")))
	prompt.WriteString(fmt.Sprintf("价格：%s\n", utils.GetStringValue(c, "price", "未提供")))
	prompt.WriteString(fmt.Sprintf("上市时间：%s\n", utils.GetStringValue(c, "releaseDate", "未提供")))
	prompt.WriteString(fmt.Sprintf("主要配置：%s\n", utils.GetStringValue(c, "specs", "")))
	prompt.WriteString(fmt.Sprintf("外观设计：%s\n", utils.GetStringValue(c, "design", "未提供")))
	prompt.WriteString(fmt.Sprintf("性能表现：%s\n", utils.GetStringValue(c, "performance", "")))
	prompt.WriteString(fmt.Sprintf("续航：%s\n", utils.GetStringValue(c, "battery", "未提供")))
	prompt.WriteString(fmt.Sprintf("相机：%s\n", utils.GetStringValue(c, "camera", "未提供")))
	prompt.WriteString(fmt.Sprintf("用户体验：%s\n", utils.GetStringValue(c, "userExperience", "")))
	prompt.WriteString(fmt.Sprintf("优点：%s\n", utils.GetStringValue(c, "pros", "")))
	prompt.WriteString(fmt.Sprintf("缺点：%s\n", utils.GetStringValue(c, "cons", "未提供")))
	prompt.WriteString(fmt.Sprintf("推荐理由：%s\n\n", utils.GetStringValue(c, "recommendation", "")))
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
