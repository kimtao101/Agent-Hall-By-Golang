
// 导入现有的 OpenAI 实例
import { openai } from './openaiService';
import logger from './logger';

interface XiaohongshuCopyRequest {
  scene: string;
  config: any;
}

interface XiaohongshuCopyResponse {
  copy: string;
}


class XiaohongshuService {
  /**
   * 生成小红书文案
   * @param request 文案请求参数
   * @returns 生成的文案
   */
  public async generateCopy(request: XiaohongshuCopyRequest): Promise<XiaohongshuCopyResponse> {
    try {
      const { scene, config } = request;
      logger.info('Generating Xiaohongshu copy', { scene, config });

      const prompt = this.getPromptForScene(scene, config);
      logger.info('Generated prompt', { prompt: prompt.substring(0, 200) + '...' });

      const response = await openai.chat.completions.create({
        messages: [
          {
            role: 'system',
            content: '你是一位专业的小红书文案撰写专家，擅长撰写各种场景的优质文案。请根据用户提供的信息，生成符合小红书平台风格的文案，包含适当的表情符号、话题标签和互动引导语。'
          },
          {
            role: 'user',
            content: prompt
          }
        ],
        model: 'deepseek-chat',
        temperature: 0.8,
        max_tokens: 1000
      });

      const copy = response.choices[0]?.message?.content || '';
      logger.info('Generated copy successfully', { copy: copy.substring(0, 200) + '...' });

      return { copy };
    } catch (error) {
      logger.error('Error generating Xiaohongshu copy', { error });
      throw error;
    }
  }

  /**
   * 根据场景生成提示词
   * @param scene 场景类型
   * @param config 场景配置
   * @returns 提示词
   */
  private getPromptForScene(scene: string, config: any): string {
    switch (scene) {
      case 'beauty':
        return this.getBeautyPrompt(config);
      case 'fashion':
        return this.getFashionPrompt(config);
      case 'travel':
        return this.getTravelPrompt(config);
      case 'food':
        return this.getFoodPrompt(config);
      case 'home':
        return this.getHomePrompt(config);
      case 'fitness':
        return this.getFitnessPrompt(config);
      case 'parenting':
        return this.getParentingPrompt(config);
      case 'tech':
        return this.getTechPrompt(config);
      default:
        return '请生成一篇小红书风格的文案';
    }
  }

  /**
   * 美妆护肤评测文案提示词
   */
  private getBeautyPrompt(config: any): string {
    const { productName, brand, price, skinType, texture, keyIngredients, usageFeel, effect, recommendation } = config;
    
    return `请为以下美妆产品生成一篇小红书风格的评测文案：

产品名称：${productName}
品牌：${brand}
价格：${price || '未提供'}
适合肤质：${skinType || '未提供'}
质地：${texture || '未提供'}
核心成分：${keyIngredients || '未提供'}
使用感受：${usageFeel}
效果：${effect}
推荐理由：${recommendation}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 产品介绍（品牌、名称、价格等基本信息）
3. 外观包装评价
4. 质地和使用感受
5. 效果评价
6. 适合人群
7. 总结推荐
8. 相关话题标签（至少5个）

要求：
- 语言风格亲切自然，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语（如"你们用过吗？"、"欢迎在评论区分享"等）
- 突出产品的核心卖点
- 内容真实可信，避免过度夸大`;
  }

  /**
   * 穿搭搭配分享文案提示词
   */
  private getFashionPrompt(config: any): string {
    const { clothingType, style, brand, price, color, material, fit, matchingTips, scenario, usageFeel } = config;
    
    return `请为以下穿搭搭配生成一篇小红书风格的分享文案：

服装类型：${clothingType}
风格：${style}
品牌：${brand || '未提供'}
价格：${price || '未提供'}
颜色：${color || '未提供'}
材质：${material || '未提供'}
版型：${fit || '未提供'}
搭配建议：${matchingTips}
适合场景：${scenario}
穿着感受：${usageFeel}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 整体搭配介绍
3. 单品推荐（材质、版型等）
4. 搭配技巧和思路
5. 适合场景和人群
6. 穿着感受
7. 购买建议
8. 相关话题标签（至少5个）

要求：
- 语言风格时尚潮流，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语
- 提供具体的搭配建议
- 突出服装的风格特点`;
  }

  /**
   * 旅行打卡攻略文案提示词
   */
  private getTravelPrompt(config: any): string {
    const { destination, duration, bestTime, budget, attractions, food, accommodation, transportation, tips, experience } = config;
    
    return `请为以下旅行目的地生成一篇小红书风格的打卡攻略文案：

目的地：${destination}
行程天数：${duration}天
最佳时间：${bestTime || '未提供'}
预算：${budget || '未提供'}
主要景点：${attractions}
特色美食：${food}
住宿推荐：${accommodation || '未提供'}
交通方式：${transportation || '未提供'}
旅行贴士：${tips}
个人体验：${experience}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 旅行概览（目的地、天数、预算等）
3. 行程安排推荐
4. 景点打卡攻略
5. 美食推荐
6. 住宿和交通建议
7. 实用贴士
8. 个人感受和总结
9. 相关话题标签（至少5个）

要求：
- 语言风格轻松愉快，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语
- 提供详细的实用信息
- 突出目的地的特色和亮点`;
  }

  /**
   * 美食探店体验文案提示词
   */
  private getFoodPrompt(config: any): string {
    const { restaurantName, location, cuisineType, priceRange, environment, service, signatureDishes, taste, recommendation } = config;
    
    return `请为以下餐厅生成一篇小红书风格的探店体验文案：

餐厅名称：${restaurantName}
位置：${location}
菜系：${cuisineType}
价格区间：${priceRange || '未提供'}
环境：${environment || '未提供'}
服务：${service || '未提供'}
招牌菜：${signatureDishes}
口味：${taste}
推荐理由：${recommendation}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 餐厅基本信息（位置、环境等）
3. 招牌菜推荐和评价
4. 口味和服务评价
5. 价格和性价比
6. 适合人群和场景
7. 总结推荐
8. 相关话题标签（至少5个）

要求：
- 语言风格生动诱人，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语
- 提供详细的菜品评价
- 突出餐厅的特色和亮点`;
  }

  /**
   * 家居好物推荐文案提示词
   */
  private getHomePrompt(config: any): string {
    const { productName, category, brand, price, material, size, usageScenario, functionality, usageFeel, spaceSaving, recommendation } = config;
    
    return `请为以下家居好物生成一篇小红书风格的推荐文案：

产品名称：${productName}
类别：${category}
品牌：${brand || '未提供'}
价格：${price || '未提供'}
材质：${material || '未提供'}
尺寸：${size || '未提供'}
使用场景：${usageScenario}
功能：${functionality}
使用感受：${usageFeel}
节省空间：${spaceSaving || '未提供'}
推荐理由：${recommendation}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 产品介绍（名称、类别、价格等）
3. 外观设计评价
4. 功能和使用方法
5. 使用感受和效果
6. 适用场景
7. 性价比评价
8. 总结推荐
9. 相关话题标签（至少5个）

要求：
- 语言风格实用亲切，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语
- 提供详细的使用体验
- 突出产品的实用价值`;
  }

  /**
   * 健身运动记录文案提示词
   */
  private getFitnessPrompt(config: any): string {
    const { workoutType, duration, frequency, equipment, difficulty, benefits, experience, tips, results } = config;
    
    return `请为以下健身运动生成一篇小红书风格的记录文案：

运动类型：${workoutType}
运动时长：${duration || '未提供'}分钟
运动频率：${frequency || '未提供'}
所需装备：${equipment || '未提供'}
难度：${difficulty || '未提供'}
运动好处：${benefits}
个人体验：${experience}
注意事项：${tips}
运动效果：${results || '未提供'}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 运动介绍（类型、时长、频率等）
3. 运动过程和感受
4. 运动好处和效果
5. 适合人群
6. 注意事项和建议
7. 个人心得和激励
8. 相关话题标签（至少5个）

要求：
- 语言风格积极向上，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语
- 提供详细的运动体验
- 传递正能量和激励信息`;
  }

  /**
   * 母婴育儿心得文案提示词
   */
  private getParentingPrompt(config: any): string {
    const { babyAge, topic, productName, brand, price, problem, solution, experience, tips } = config;
    
    return `请为以下母婴育儿主题生成一篇小红书风格的心得分享文案：

宝宝年龄：${babyAge}
主题：${topic}
产品名称：${productName || '未提供'}
品牌：${brand || '未提供'}
价格：${price || '未提供'}
问题描述：${problem}
解决方案：${solution}
经验分享：${experience}
小贴士：${tips}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 问题背景介绍
3. 解决方案分享
4. 经验总结
5. 实用小贴士
6. 产品推荐（如果有）
7. 互动和鼓励
8. 相关话题标签（至少5个）

要求：
- 语言风格温暖亲切，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语
- 提供实用的育儿经验
- 传递正能量和支持信息`;
  }

  /**
   * 数码产品测评文案提示词
   */
  private getTechPrompt(config: any): string {
    const { productName, brand, price, releaseDate, specs, design, performance, battery, camera, userExperience, pros, cons, recommendation } = config;
    
    return `请为以下数码产品生成一篇小红书风格的测评文案：

产品名称：${productName}
品牌：${brand}
价格：${price || '未提供'}
上市时间：${releaseDate || '未提供'}
主要配置：${specs}
外观设计：${design || '未提供'}
性能表现：${performance}
续航：${battery || '未提供'}
相机：${camera || '未提供'}
用户体验：${userExperience}
优点：${pros}
缺点：${cons || '未提供'}
推荐理由：${recommendation}

请按照以下结构撰写：
1. 吸引人的标题（包含emoji）
2. 产品开箱介绍
3. 外观设计评价
4. 性能和功能测试
5. 用户体验
6. 优缺点分析
7. 适合人群
8. 性价比评价
9. 总结推荐
10. 相关话题标签（至少5个）

要求：
- 语言风格专业易懂，符合小红书用户的阅读习惯
- 适当使用emoji表情符号
- 包含互动引导语
- 提供详细的产品信息
- 评价客观公正，避免过度 bias`;
  }
}

// 导出单例实例
export const xiaohongshuService = new XiaohongshuService();
export { XiaohongshuService };
