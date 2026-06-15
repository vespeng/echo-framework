package user

import (
	"echo-framework/internal/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

var svc = NewService()

// RegisterRoutes 注册用户模块路由
func RegisterRoutes(g *echo.Group) {
	g.GET("", list)
	g.GET("/:id", get)
	g.POST("", create)
	g.PUT("/:id", update)
	g.DELETE("/:id", remove)
}

// list 获取用户列表
func list(c echo.Context) error {
	users, err := svc.GetAllUsers()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    500,
			"message": "获取用户列表失败",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"message": "success",
		"data":    users,
	})
}

// get 根据ID获取用户
func get(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"message": "无效的用户ID",
		})
	}
	u, err := svc.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"code":    404,
			"message": "用户不存在",
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"message": "success",
		"data":    u,
	})
}

// create 创建用户
func create(c echo.Context) error {
	var u model.SysUser
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"message": "请求参数错误",
		})
	}
	if err := svc.AddUser(&u); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    500,
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusCreated, echo.Map{
		"code":    201,
		"message": "创建成功",
		"data":    u,
	})
}

// update 更新用户
func update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"message": "无效的用户ID",
		})
	}
	var u model.SysUser
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"message": "请求参数错误",
		})
	}
	if err := svc.ModifyUser(id, &u); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    500,
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"message": "更新成功",
	})
}

// remove 删除用户
func remove(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"code":    400,
			"message": "无效的用户ID",
		})
	}
	if err := svc.RemoveUser(id); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"code":    500,
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"code":    200,
		"message": "删除成功",
	})
}
