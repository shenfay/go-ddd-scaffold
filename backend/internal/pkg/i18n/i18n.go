package i18n

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
)

// Language 支持的语言
type Language string

const (
	LangZHCN Language = "zh-CN" // 简体中文
	LangENUS Language = "en-US" // 英语(美国)
	LangJAJP Language = "ja-JP" // 日语
)

// 默认语言
const DefaultLang = LangZHCN

// ContextKey 语言上下文键
type ContextKey string

const LangContextKey ContextKey = "lang"

// messages 消息映射表
var messages = map[Language]map[string]string{
	LangZHCN: {
		// 通用消息
		"success":           "操作成功",
		"failed":            "操作失败",
		"system_error":      "系统内部错误，请稍后重试",
		"validation_failed": "参数校验失败",
		"not_found":         "请求的资源不存在",
		"unauthorized":      "未授权，请先登录",
		"forbidden":         "没有权限访问该资源",

		// 知识领域相关
		"domain_not_found":     "知识领域不存在",
		"domain_create_failed": "创建知识领域失败",
		"domain_update_failed": "更新知识领域失败",
		"domain_delete_failed": "删除知识领域失败",

		// 主干相关
		"trunk_not_found":     "主干不存在",
		"trunk_create_failed": "创建主干失败",

		// 节点相关
		"node_not_found":     "节点不存在",
		"node_create_failed": "创建节点失败",
		"node_update_failed": "更新节点失败",
		"node_delete_failed": "删除节点失败",

		// 关系相关
		"relationship_exists":    "关系已存在",
		"relationship_not_found": "关系不存在",
		"self_reference_error":   "节点不能与自己建立关系",

		// 分页相关
		"page_param_error": "分页参数错误",

		// 知识领域翻译
		"domains.numeracy_algebra.name":        "数与代数",
		"domains.numeracy_algebra.description": "学习数与代数基础知识，包括数的认识、运算、方程等",
		"domains.geometry.name":                "图形与几何",
		"domains.geometry.description":         "认识基本几何图形，学习测量与空间方位",
		"domains.measurement.name":             "计量",
		"domains.measurement.description":      "学习时间、长度、重量等计量知识",
		"domains.statistics.name":              "统计与概率",
		"domains.statistics.description":       "了解数据收集、整理与概率基础知识",

		// 主干翻译
		"trunk.addition.name":               "加法",
		"trunk.addition.description":        "整数的加法运算",
		"trunk.subtraction.name":            "减法",
		"trunk.subtraction.description":     "整数的减法运算",
		"trunk.multiplication.name":         "乘法",
		"trunk.multiplication.description":  "整数的乘法运算",
		"trunk.division.name":               "除法",
		"trunk.division.description":        "整数的除法运算",
		"trunk.fraction.name":               "分数",
		"trunk.fraction.description":        "分数的认识与运算",
		"trunk.decimal.name":                "小数",
		"trunk.decimal.description":         "小数的认识与运算",
		"trunk.basic_shapes.name":           "基本图形",
		"trunk.basic_shapes.description":    "认识基本几何图形",
		"trunk.measurement.name":            "测量",
		"trunk.measurement.description":     "长度、面积、体积的测量",
		"trunk.spatial.name":                "空间与方位",
		"trunk.spatial.description":         "认识空间方位和位置",
		"trunk.time.name":                   "时间",
		"trunk.time.description":            "认识时钟和时间计算",
		"trunk.money.name":                  "货币",
		"trunk.money.description":           "认识货币和简单计算",
		"trunk.length.name":                 "长度",
		"trunk.length.description":          "长度的认识与测量",
		"trunk.data_collection.name":        "数据收集",
		"trunk.data_collection.description": "数据的收集与整理",
		"trunk.probability.name":            "概率",
		"trunk.probability.description":     "概率基础知识",

		// 节点翻译
		"nodes.addition_concept_lv1.name_child":  "什么是加法",
		"nodes.addition_concept_lv1.name_parent": "加法的意义",
		"nodes.addition_concept_lv1.description": "理解加法是把两部分合在一起的概念",
		"nodes.addition_concept_lv2.name_child":  "加法的变化",
		"nodes.addition_concept_lv2.name_parent": "加法各部分关系",
		"nodes.addition_concept_lv2.description": "理解加数、被加数、和的关系",
		"nodes.addition_concept_lv3.name_child":  "复杂的加法",
		"nodes.addition_concept_lv3.name_parent": "多位数加法",
		"nodes.addition_concept_lv3.description": "掌握多位数加法的计算方法",

		"nodes.addition_skill_lv1.name_child":  "学做加法",
		"nodes.addition_skill_lv1.name_parent": "10以内加法",
		"nodes.addition_skill_lv1.description": "掌握10以内数的加法计算",
		"nodes.addition_skill_lv2.name_child":  "更大数的加法",
		"nodes.addition_skill_lv2.name_parent": "20以内加法",
		"nodes.addition_skill_lv2.description": "掌握20以内数的加法计算",
		"nodes.addition_skill_lv3.name_child":  "加减混合",
		"nodes.addition_skill_lv3.name_parent": "加减混合运算",
		"nodes.addition_skill_lv3.description": "掌握加减混合运算顺序",

		"nodes.fraction_concept_lv2.name_child":  "认识分数",
		"nodes.fraction_concept_lv2.name_parent": "分数的意义",
		"nodes.fraction_concept_lv2.description": "理解分数表示整体的一部分",
		"nodes.fraction_concept_lv3.name_child":  "分数大小",
		"nodes.fraction_concept_lv3.name_parent": "分数的大小比较",
		"nodes.fraction_concept_lv3.description": "学会比较分数的大小",
		"nodes.fraction_skill_lv2.name_child":    "分数计算入门",
		"nodes.fraction_skill_lv2.name_parent":   "同分母分数加减",
		"nodes.fraction_skill_lv2.description":   "掌握同分母分数加减法",
		"nodes.fraction_skill_lv3.name_child":    "分数进阶",
		"nodes.fraction_skill_lv3.name_parent":   "异分母分数加减",
		"nodes.fraction_skill_lv3.description":   "掌握异分母分数加减法",
		"nodes.fraction_problem_lv4.name_child":  "分数应用",
		"nodes.fraction_problem_lv4.name_parent": "分数应用题",
		"nodes.fraction_problem_lv4.description": "解决分数相关的实际问题",

		"nodes.circle_concept_lv1.name_child":    "圆形朋友",
		"nodes.circle_concept_lv1.name_parent":   "圆的认识",
		"nodes.circle_concept_lv1.description":   "认识圆的基本特征",
		"nodes.triangle_concept_lv2.name_child":  "三角形",
		"nodes.triangle_concept_lv2.name_parent": "三角形的认识",
		"nodes.triangle_concept_lv2.description": "认识三角形及其特性",

		"nodes.time_concept_lv1.name_child":  "认识时间",
		"nodes.time_concept_lv1.name_parent": "时间概念",
		"nodes.time_concept_lv1.description": "认识时钟和时间",
		"nodes.time_skill_lv1.name_child":    "看时间",
		"nodes.time_skill_lv1.name_parent":   "读取时钟时间",
		"nodes.time_skill_lv1.description":   "学会看时钟读取时间",

		"nodes.probability_concept_lv3.name_child":  "可能性",
		"nodes.probability_concept_lv3.name_parent": "可能性大小",
		"nodes.probability_concept_lv3.description": "理解可能性的大小",
		"nodes.probability_skill_lv4.name_child":    "算概率",
		"nodes.probability_skill_lv4.name_parent":   "简单概率计算",
		"nodes.probability_skill_lv4.description":   "学会计算简单概率",

		// 标签翻译
		"tags.fraction.name":           "分数",
		"tags.fraction.description":    "与分数相关的知识点",
		"tags.addition.name":           "加法",
		"tags.addition.description":    "与加法相关的知识点",
		"tags.time.name":               "时间",
		"tags.time.description":        "与时间相关的知识点",
		"tags.probability.name":        "概率",
		"tags.probability.description": "与概率相关的知识点",

		// 能力等级翻译
		"levels.bloom_solo_fusion.lv1.description": "感知级：通过观察和具体例子理解概念",
		"levels.bloom_solo_fusion.lv2.description": "操作级：能够进行基本的操作和计算",
		"levels.bloom_solo_fusion.lv3.description": "规则级：理解规则并能运用规则解决问题",
		"levels.bloom_solo_fusion.lv4.description": "应用级：将知识应用于新的情境",
		"levels.bloom_solo_fusion.lv5.description": "创新级：能够创造性地解决问题",
	},
	LangENUS: {
		// 通用消息
		"success":           "Operation successful",
		"failed":            "Operation failed",
		"system_error":      "Internal system error, please try again later",
		"validation_failed": "Parameter validation failed",
		"not_found":         "Resource not found",
		"unauthorized":      "Unauthorized, please login first",
		"forbidden":         "Access denied",

		// 知识领域相关
		"domain_not_found":     "Knowledge domain does not exist",
		"domain_create_failed": "Failed to create knowledge domain",
		"domain_update_failed": "Failed to update knowledge domain",
		"domain_delete_failed": "Failed to delete knowledge domain",

		// 主干相关
		"trunk_not_found":     "Trunk does not exist",
		"trunk_create_failed": "Failed to create trunk",

		// 节点相关
		"node_not_found":     "Node does not exist",
		"node_create_failed": "Failed to create node",
		"node_update_failed": "Failed to update node",
		"node_delete_failed": "Failed to delete node",

		// 关系相关
		"relationship_exists":    "Relationship already exists",
		"relationship_not_found": "Relationship not found",
		"self_reference_error":   "Node cannot establish relationship with itself",

		// 分页相关
		"page_param_error": "Pagination parameter error",

		// 知识领域翻译
		"domains.numeracy_algebra.name":        "Numbers and Algebra",
		"domains.numeracy_algebra.description": "Learn basic algebraic knowledge including number recognition, operations, equations, etc.",
		"domains.geometry.name":                "Geometry",
		"domains.geometry.description":         "Learn basic geometric shapes, measurement and spatial orientation",
		"domains.measurement.name":             "Measurement",
		"domains.measurement.description":      "Learn measurement knowledge of time, length, weight, etc.",
		"domains.statistics.name":              "Statistics and Probability",
		"domains.statistics.description":       "Learn basic data collection, organization and probability",

		// 主干翻译
		"trunk.addition.name":               "Addition",
		"trunk.addition.description":        "Integer addition operations",
		"trunk.subtraction.name":            "Subtraction",
		"trunk.subtraction.description":     "Integer subtraction operations",
		"trunk.multiplication.name":         "Multiplication",
		"trunk.multiplication.description":  "Integer multiplication operations",
		"trunk.division.name":               "Division",
		"trunk.division.description":        "Integer division operations",
		"trunk.fraction.name":               "Fractions",
		"trunk.fraction.description":        "Understanding and operations of fractions",
		"trunk.decimal.name":                "Decimals",
		"trunk.decimal.description":         "Understanding and operations of decimals",
		"trunk.basic_shapes.name":           "Basic Shapes",
		"trunk.basic_shapes.description":    "Learn basic geometric shapes",
		"trunk.measurement.name":            "Measurement",
		"trunk.measurement.description":     "Measurement of length, area, and volume",
		"trunk.spatial.name":                "Space and Position",
		"trunk.spatial.description":         "Understanding spatial orientation and positions",
		"trunk.time.name":                   "Time",
		"trunk.time.description":            "Telling time and time calculations",
		"trunk.money.name":                  "Money",
		"trunk.money.description":           "Understanding money and simple calculations",
		"trunk.length.name":                 "Length",
		"trunk.length.description":          "Understanding and measuring length",
		"trunk.data_collection.name":        "Data Collection",
		"trunk.data_collection.description": "Data collection and organization",
		"trunk.probability.name":            "Probability",
		"trunk.probability.description":     "Basic probability concepts",

		// 节点翻译
		"nodes.addition_concept_lv1.name_child":  "What is Addition",
		"nodes.addition_concept_lv1.name_parent": "Meaning of Addition",
		"nodes.addition_concept_lv1.description": "Understanding addition as combining two parts",
		"nodes.addition_concept_lv2.name_child":  "Changes in Addition",
		"nodes.addition_concept_lv2.name_parent": "Parts of Addition",
		"nodes.addition_concept_lv2.description": "Understanding addend, sum relationships",
		"nodes.addition_concept_lv3.name_child":  "Complex Addition",
		"nodes.addition_concept_lv3.name_parent": "Multi-digit Addition",
		"nodes.addition_concept_lv3.description": "Master multi-digit addition",

		"nodes.addition_skill_lv1.name_child":  "Learn Addition",
		"nodes.addition_skill_lv1.name_parent": "Addition within 10",
		"nodes.addition_skill_lv1.description": "Master addition within 10",
		"nodes.addition_skill_lv2.name_child":  "Larger Numbers",
		"nodes.addition_skill_lv2.name_parent": "Addition within 20",
		"nodes.addition_skill_lv2.description": "Master addition within 20",
		"nodes.addition_skill_lv3.name_child":  "Mixed Operations",
		"nodes.addition_skill_lv3.name_parent": "Mixed Addition and Subtraction",
		"nodes.addition_skill_lv3.description": "Master order of operations",

		"nodes.fraction_concept_lv2.name_child":  "Understanding Fractions",
		"nodes.fraction_concept_lv2.name_parent": "Meaning of Fractions",
		"nodes.fraction_concept_lv2.description": "Fractions represent parts of a whole",
		"nodes.fraction_concept_lv3.name_child":  "Fraction Size",
		"nodes.fraction_concept_lv3.name_parent": "Comparing Fractions",
		"nodes.fraction_concept_lv3.description": "Learn to compare fractions",
		"nodes.fraction_skill_lv2.name_child":    "Fraction Basics",
		"nodes.fraction_skill_lv2.name_parent":   "Same Denominator Fractions",
		"nodes.fraction_skill_lv2.description":   "Add/subtract fractions with same denominator",
		"nodes.fraction_skill_lv3.name_child":    "Fraction Advanced",
		"nodes.fraction_skill_lv3.name_parent":   "Different Denominator Fractions",
		"nodes.fraction_skill_lv3.description":   "Add/subtract fractions with different denominators",
		"nodes.fraction_problem_lv4.name_child":  "Fraction Applications",
		"nodes.fraction_problem_lv4.name_parent": "Fraction Word Problems",
		"nodes.fraction_problem_lv4.description": "Solve real-world fraction problems",

		"nodes.circle_concept_lv1.name_child":    "Circle Friend",
		"nodes.circle_concept_lv1.name_parent":   "Understanding Circles",
		"nodes.circle_concept_lv1.description":   "Learn basic properties of circles",
		"nodes.triangle_concept_lv2.name_child":  "Triangles",
		"nodes.triangle_concept_lv2.name_parent": "Understanding Triangles",
		"nodes.triangle_concept_lv2.description": "Learn about triangles and their properties",

		"nodes.time_concept_lv1.name_child":  "Understanding Time",
		"nodes.time_concept_lv1.name_parent": "Time Concepts",
		"nodes.time_concept_lv1.description": "Learn about clocks and time",
		"nodes.time_skill_lv1.name_child":    "Telling Time",
		"nodes.time_skill_lv1.name_parent":   "Reading Clock Time",
		"nodes.time_skill_lv1.description":   "Learn to read time on clocks",

		"nodes.probability_concept_lv3.name_child":  "Possibility",
		"nodes.probability_concept_lv3.name_parent": "Probability",
		"nodes.probability_concept_lv3.description": "Understanding likelihood",
		"nodes.probability_skill_lv4.name_child":    "Calculating Probability",
		"nodes.probability_skill_lv4.name_parent":   "Simple Probability",
		"nodes.probability_skill_lv4.description":   "Calculate simple probability",

		// 标签翻译
		"tags.fraction.name":           "Fractions",
		"tags.fraction.description":    "Knowledge points related to fractions",
		"tags.addition.name":           "Addition",
		"tags.addition.description":    "Knowledge points related to addition",
		"tags.time.name":               "Time",
		"tags.time.description":        "Knowledge points related to time",
		"tags.probability.name":        "Probability",
		"tags.probability.description": "Knowledge points related to probability",

		// 能力等级翻译
		"levels.bloom_solo_fusion.lv1.description": "Perception: Understand concepts through observation and examples",
		"levels.bloom_solo_fusion.lv2.description": "Operation: Perform basic operations and calculations",
		"levels.bloom_solo_fusion.lv3.description": "Rule: Understand and apply rules to solve problems",
		"levels.bloom_solo_fusion.lv4.description": "Application: Apply knowledge to new situations",
		"levels.bloom_solo_fusion.lv5.description": "Innovation: Solve problems creatively",
	},
	LangJAJP: {
		// 通用消息
		"success":           "操作成功",
		"failed":            "操作失敗",
		"system_error":      "システムエラーが発生しました。もう一度お試しください",
		"validation_failed": "パラメータ検証に失敗しました",
		"not_found":         "リソースが見つかりません",
		"unauthorized":      "認証されていません。ログインしてください",
		"forbidden":         "アクセスが拒否されました",

		// 知识领域相关
		"domain_not_found":     "ナレッジドメインが存在しません",
		"domain_create_failed": "ナレッジドメインの作成に失敗しました",
		"domain_update_failed": "ナレッジドメインの更新に失敗しました",
		"domain_delete_failed": "ナレッジドメインの削除に失敗しました",

		// 主干相关
		"trunk_not_found":     "トランクが存在しません",
		"trunk_create_failed": "トランクの作成に失敗しました",

		// 节点相关
		"node_not_found":     "ノードが存在しません",
		"node_create_failed": "ノードの作成に失敗しました",
		"node_update_failed": "ノードの更新に失敗しました",
		"node_delete_failed": "ノードの削除に失敗しました",

		// 关系相关
		"relationship_exists":    "関係は既に存在します",
		"relationship_not_found": "関係が見つかりません",
		"self_reference_error":   "ノードは自分自身との関係を設定できません",

		// 分页相关
		"page_param_error": "ページングパラメータエラー",

		// 知识领域翻译
		"domains.numeracy_algebra.name":        "数と代数",
		"domains.numeracy_algebra.description": "数の認識、演算、方程式などの代数学の基礎を学ぶ",
		"domains.geometry.name":                "図形と幾何",
		"domains.geometry.description":         "基本的な図形を学び、測定と空間方位を学ぶ",
		"domains.measurement.name":             "计量",
		"domains.measurement.description":      "時間、長さ、重さなどの計量の知識を学ぶ",
		"domains.statistics.name":              "統計と確率",
		"domains.statistics.description":       "データ収集、整理と確率の基礎を学ぶ",

		// 主干翻译
		"trunk.addition.name":               "足し算",
		"trunk.addition.description":        "整数の足し算",
		"trunk.subtraction.name":            "引き算",
		"trunk.subtraction.description":     "整数の引き算",
		"trunk.multiplication.name":         "掛け算",
		"trunk.multiplication.description":  "整数の掛け算",
		"trunk.division.name":               "割り算",
		"trunk.division.description":        "整数の割り算",
		"trunk.fraction.name":               "分数",
		"trunk.fraction.description":        "分数の理解と計算",
		"trunk.decimal.name":                "小数",
		"trunk.decimal.description":         "小数の理解と計算",
		"trunk.basic_shapes.name":           "基本図形",
		"trunk.basic_shapes.description":    "基本幾何図形を学ぶ",
		"trunk.measurement.name":            "測定",
		"trunk.measurement.description":     "長さ、面積、体積の測定",
		"trunk.spatial.name":                "空間と位置",
		"trunk.spatial.description":         "空間方位と位置の理解",
		"trunk.time.name":                   "時間",
		"trunk.time.description":            "時計の読みと時間計算",
		"trunk.money.name":                  "お金",
		"trunk.money.description":           "お金と簡単な計算",
		"trunk.length.name":                 "長さ",
		"trunk.length.description":          "長さの理解と測定",
		"trunk.data_collection.name":        "データ収集",
		"trunk.data_collection.description": "データの収集と整理",
		"trunk.probability.name":            "確率",
		"trunk.probability.description":     "確率の基礎",

		// 节点翻译
		"nodes.addition_concept_lv1.name_child":  "足し算ってなに？",
		"nodes.addition_concept_lv1.name_parent": "足し算の意味",
		"nodes.addition_concept_lv1.description": "足し算は两部分を合わせることです",
		"nodes.addition_concept_lv2.name_child":  "足し算の変化",
		"nodes.addition_concept_lv2.name_parent": "足し算の関係",
		"nodes.addition_concept_lv2.description": "加数、被加数、和の関係",
		"nodes.addition_concept_lv3.name_child":  "複雑な足し算",
		"nodes.addition_concept_lv3.name_parent": "多桁の足し算",
		"nodes.addition_concept_lv3.description": "多桁の足し算をマスター",

		"nodes.addition_skill_lv1.name_child":  "足し算を学ぼう",
		"nodes.addition_skill_lv1.name_parent": "10以内の足し算",
		"nodes.addition_skill_lv1.description": "10以内の足し算をマスター",
		"nodes.addition_skill_lv2.name_child":  "もっと大きい数",
		"nodes.addition_skill_lv2.name_parent": "20以内の足し算",
		"nodes.addition_skill_lv2.description": "20以内の足し算をマスター",
		"nodes.addition_skill_lv3.name_child":  "混合計算",
		"nodes.addition_skill_lv3.name_parent": "足し算と引き算",
		"nodes.addition_skill_lv3.description": "計算の順序をマスター",

		"nodes.fraction_concept_lv2.name_child":  "分数を調べよう",
		"nodes.fraction_concept_lv2.name_parent": "分数の意味",
		"nodes.fraction_concept_lv2.description": "分数は全体の一部を表す",
		"nodes.fraction_concept_lv3.name_child":  "分数の大きさ",
		"nodes.fraction_concept_lv3.name_parent": "分数の比較",
		"nodes.fraction_concept_lv3.description": "分数の大きさを比べられる",
		"nodes.fraction_skill_lv2.name_child":    "分数の計算",
		"nodes.fraction_skill_lv2.name_parent":   "同じ分母の分数",
		"nodes.fraction_skill_lv2.description":   "分母が同じ分数の足し算と引き算",
		"nodes.fraction_skill_lv3.name_child":    "分数の上級",
		"nodes.fraction_skill_lv3.name_parent":   "異なる分母の分数",
		"nodes.fraction_skill_lv3.description":   "分母が異なる分数の足し算と引き算",
		"nodes.fraction_problem_lv4.name_child":  "分数の問題",
		"nodes.fraction_problem_lv4.name_parent": "分数文章問題",
		"nodes.fraction_problem_lv4.description": "実際の分数の問題を解く",

		"nodes.circle_concept_lv1.name_child":    "まる",
		"nodes.circle_concept_lv1.name_parent":   "円の理解",
		"nodes.circle_concept_lv1.description":   "円の基本的な性質",
		"nodes.triangle_concept_lv2.name_child":  "三角形",
		"nodes.triangle_concept_lv2.name_parent": "三角形の理解",
		"nodes.triangle_concept_lv2.description": "三角形とその特性",

		"nodes.time_concept_lv1.name_child":  "時間を学ぼう",
		"nodes.time_concept_lv1.name_parent": "時間の概念",
		"nodes.time_concept_lv1.description": "時計と時間を学ぶ",
		"nodes.time_skill_lv1.name_child":    "時刻を見る",
		"nodes.time_skill_lv1.name_parent":   "時計の読み方",
		"nodes.time_skill_lv1.description":   "時計を読む",

		"nodes.probability_concept_lv3.name_child":  "可能性",
		"nodes.probability_concept_lv3.name_parent": "確率",
		"nodes.probability_concept_lv3.description": "可能性大小を理解",
		"nodes.probability_skill_lv4.name_child":    "確率を計算",
		"nodes.probability_skill_lv4.name_parent":   "簡単な確率",
		"nodes.probability_skill_lv4.description":   "簡単な確率を計算",

		// 标签翻译
		"tags.fraction.name":           "分数",
		"tags.fraction.description":    "分数関連の知识点",
		"tags.addition.name":           "足し算",
		"tags.addition.description":    "足し算関連の知识点",
		"tags.time.name":               "時間",
		"tags.time.description":        "時間関連の知识点",
		"tags.probability.name":        "確率",
		"tags.probability.description": "確率関連の知识点",

		// 能力等级翻译
		"levels.bloom_solo_fusion.lv1.description": "感知級：観察と例を通じて概念を理解",
		"levels.bloom_solo_fusion.lv2.description": "操作級：基本的な操作と計算ができる",
		"levels.bloom_solo_fusion.lv3.description": "規則級：規則を理解し問題解決に適用",
		"levels.bloom_solo_fusion.lv4.description": "応用級：知識を新しい状況に適用",
		"levels.bloom_solo_fusion.lv5.description": "創造級：創造的に問題解決",
	},
}

var (
	mu          sync.RWMutex
	defaultLang = DefaultLang
	// 语言文件加载器
	localeLoader = GetDefaultLoader()
	localeOnce   sync.Once
)

// initLocale 初始化语言文件
func initLocale() {
	localeOnce.Do(func() {
		// 尝试加载locale目录下的YAML文件
		// 生产环境可配置实际路径
		localePath := "./internal/pkg/i18n/locale"
		_ = GetDefaultLoader().LoadFromDirectory(localePath)
	})
}

// LoadLocaleFiles 手动加载语言文件（可在main中调用）
func LoadLocaleFiles(dirPath string) error {
	return GetDefaultLoader().LoadFromDirectory(dirPath)
}

// SetDefaultLang 设置默认语言
func SetDefaultLang(lang Language) {
	mu.Lock()
	defer mu.Unlock()
	defaultLang = lang
}

// GetDefaultLang 获取默认语言
func GetDefaultLang() Language {
	mu.RLock()
	defer mu.RUnlock()
	return defaultLang
}

// GetMessage 获取消息
func GetMessage(lang Language, key string) string {
	// 确保语言文件已加载
	initLocale()

	// 1. 优先从YAML语言文件获取
	if msg, ok := GetDefaultLoader().Get(string(lang), key); ok {
		return msg
	}

	mu.RLock()
	defer mu.RUnlock()

	// 2. 回退到内置消息映射
	if msgMap, ok := messages[lang]; ok {
		if msg, ok := msgMap[key]; ok {
			return msg
		}
	}

	// 3. 回退到默认语言
	if msgMap, ok := messages[defaultLang]; ok {
		if msg, ok := msgMap[key]; ok {
			return msg
		}
	}

	// 返回key作为最后回退
	return key
}

// GetMessageByContext 从上下文获取语言并返回消息
func GetMessageByContext(ctx context.Context, key string) string {
	lang := GetLangFromContext(ctx)
	return GetMessage(lang, key)
}

// GetLangFromContext 从上下文获取语言
func GetLangFromContext(ctx context.Context) Language {
	if lang, ok := ctx.Value(LangContextKey).(Language); ok {
		return lang
	}
	return GetDefaultLang()
}

// SetLangToContext 将语言设置到上下文
func SetLangToContext(ctx context.Context, lang Language) context.Context {
	return context.WithValue(ctx, LangContextKey, lang)
}

// GetLangFromGinContext 从Gin上下文获取语言
func GetLangFromGinContext(c *gin.Context) Language {
	// 1. 优先从URL参数获取
	if lang := c.Query("lang"); lang != "" {
		if l := parseLanguage(lang); l != "" {
			return l
		}
	}

	// 2. 从Header获取
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang != "" {
		if l := parseAcceptLanguage(acceptLang); l != "" {
			return l
		}
	}

	// 3. 从Cookie获取
	if lang, err := c.Cookie("lang"); err == nil {
		if l := parseLanguage(lang); l != "" {
			return l
		}
	}

	return GetDefaultLang()
}

// parseLanguage 解析语言字符串
func parseLanguage(lang string) Language {
	switch lang {
	case "zh-CN", "zh", "CN":
		return LangZHCN
	case "en-US", "en", "US":
		return LangENUS
	case "ja-JP", "ja", "JP":
		return LangJAJP
	default:
		return ""
	}
}

// parseAcceptLanguage 解析Accept-Language头
func parseAcceptLanguage(acceptLang string) Language {
	// 简单处理：取第一个语言
	for _, lang := range []string{"zh-CN", "zh", "en-US", "en", "ja-JP", "ja"} {
		if len(acceptLang) >= len(lang) && acceptLang[:len(lang)] == lang {
			return parseLanguage(lang)
		}
	}
	return ""
}

// AddMessage 添加自定义消息
func AddMessage(lang Language, key, msg string) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := messages[lang]; !ok {
		messages[lang] = make(map[string]string)
	}
	messages[lang][key] = msg
}

// AddMessages 批量添加消息
func AddMessages(lang Language, msgs map[string]string) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := messages[lang]; !ok {
		messages[lang] = make(map[string]string)
	}

	for k, v := range msgs {
		messages[lang][k] = v
	}
}
