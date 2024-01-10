package v1

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"varp_be/internal/entity"
	"varp_be/pkg/logger"

	"github.com/gin-gonic/gin"
	goonvif "github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
)

type cameraRoutes struct {
	l logger.Interface
}

func newCamneraRoutes(handler *gin.RouterGroup, l logger.Interface) {
	r := &cameraRoutes{l}

	h := handler.Group("/camera")
	{
		h.GET("/device-on-local", r.loadCameras)
	}
}

// @Summary     Translate
// @Description detect camera in LAN network
// @ID          do-translate
// @Tags  	    translation
// @Accept      json
// @Produce     json
// @Param       request body doTranslateRequest true "Set up translation"
// @Success     200 {object} entity.Translation
// @Failure     400 {object} response
// @Failure     500 {object} response
// @Router      /translation/do-translate [post]
func (r *cameraRoutes) loadCameras(c *gin.Context) {
	// Set the timeout for the discovery process

	// Perform ONVIF discovery
	ethernetInterface := getNetworkInterfaceName()
	devices, err := goonvif.GetAvailableDevicesAtSpecificEthernetInterface(ethernetInterface)
	if err != nil {
		log.Fatal("Failed to perform ONVIF discovery:", err)
	}

	cameras := make([]entity.Camera, 0)
	for _, dev := range devices {

		getCapabilities := device.GetCapabilities{Category: "All"}
		capabilities, err := dev.CallMethod(getCapabilities)

		if err != nil {
			log.Fatal("Failed to get LAN IP address:", err, capabilities)
		}

		cameras = append(cameras, entity.Camera{})
	}

	c.JSON(http.StatusOK, devices)
}

func getNetworkInterfaceName() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Lỗi khi lấy danh sách network interface:", err)
		return ""
	}
	for _, iface := range interfaces {
		if (iface.Flags & net.FlagUp) != 0 {
			return iface.Name
		}
	}
	return ""
}

func getLANIPAddress(xaddrs []string) (string, error) {
	for _, addr := range xaddrs {
		host, _, err := net.SplitHostPort(addr)
		if err == nil {
			ip := net.ParseIP(host)
			if ip != nil && !ip.IsLoopback() && ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("LAN IP address not found")
}
