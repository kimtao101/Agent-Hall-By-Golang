// xiaohongshu 提供小红书文案生成智能体功能
package xiaohongshu

// CopyRequest 文案生成请求
type CopyRequest struct {
	Scene  string                 `json:"scene" binding:"required"`
	Config map[string]interface{} `json:"config" binding:"required"`
}

// CopyResponse 文案生成响应
type CopyResponse struct {
	Copy string `json:"copy"`
}

// Scene 定义文案场景
type Scene struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

// AvailableScenes 返回所有可用的文案场景
func AvailableScenes() []Scene {
	return []Scene{
		{
			ID:          "beauty",
			Name:        "美妆护肤评测",
			Icon:        "💄",
			Description: "产品评测、使用感受、效果分享",
		},
		{
			ID:          "fashion",
			Name:        "穿搭搭配分享",
			Icon:        "👗",
			Description: "服装搭配、风格推荐、购物指南",
		},
		{
			ID:          "travel",
			Name:        "旅行打卡攻略",
			Icon:        "✈️",
			Description: "景点推荐、行程规划、旅行体验",
		},
		{
			ID:          "food",
			Name:        "美食探店体验",
			Icon:        "🍔",
			Description: "餐厅评测、美食推荐、用餐体验",
		},
		{
			ID:          "home",
			Name:        "家居好物推荐",
			Icon:        "🏠",
			Description: "家居用品、装修灵感、生活技巧",
		},
		{
			ID:          "fitness",
			Name:        "健身运动记录",
			Icon:        "🏋️",
			Description: "运动计划、健身心得、成果分享",
		},
		{
			ID:          "parenting",
			Name:        "母婴育儿心得",
			Icon:        "👶",
			Description: "育儿经验、产品推荐、成长记录",
		},
		{
			ID:          "tech",
			Name:        "数码产品测评",
			Icon:        "📱",
			Description: "产品评测、使用体验、技术分析",
		},
	}
}
