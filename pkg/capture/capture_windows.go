package capture

import (
	"image"
	"syscall"
	"unsafe"
)

const DX_DLL = "d3d9.dll"

const (
	DEVTYPE_HAL                      = 1
	CREATE_SOFTWARE_VERTEXPROCESSING = 0x00000020
	SWAPEFFECT_DISCARD               = 1
	FMT_A8R8G8B8                     = 21
	POOL_SYSTEMMEM                   = 2
)

type MULTISAMPLE_TYPE uint32
type SWAPEFFECT uint32
type POOL uint32
type DEVTYPE uint32
type FORMAT uint32

type Direct3D struct {
	vtbl *direct3DVtbl
}

type direct3DVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	RegisterSoftwareDevice      uintptr
	GetAdapterCount             uintptr
	GetAdapterIdentifier        uintptr
	GetAdapterModeCount         uintptr
	EnumAdapterModes            uintptr
	GetAdapterDisplayMode       uintptr
	CheckDeviceType             uintptr
	CheckDeviceFormat           uintptr
	CheckDeviceMultiSampleType  uintptr
	CheckDepthStencilMatch      uintptr
	CheckDeviceFormatConversion uintptr
	GetDeviceCaps               uintptr
	GetAdapterMonitor           uintptr
	CreateDevice                uintptr
}

type DISPLAYMODE struct {
	Width       uint32
	Height      uint32
	RefreshRate uint32
	Format      FORMAT
}

func (obj *Direct3D) Release() uint32 {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return uint32(ret)
}

func (obj *Direct3D) GetAdapterDisplayMode(adapter uint) (mode DISPLAYMODE, err error) {
	syscall.Syscall(
		obj.vtbl.GetAdapterDisplayMode,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(adapter),
		uintptr(unsafe.Pointer(&mode)),
	)
	return
}

type (
	HANDLE   uintptr
	HWND     HANDLE
	HMONITOR HANDLE
	HDC      HANDLE
)

type PRESENT_PARAMETERS struct {
	BackBufferWidth            uint32
	BackBufferHeight           uint32
	BackBufferFormat           FORMAT
	BackBufferCount            uint32
	MultiSampleType            MULTISAMPLE_TYPE
	MultiSampleQuality         uint32
	SwapEffect                 SWAPEFFECT
	HDeviceWindow              HWND
	Windowed                   int32
	EnableAutoDepthStencil     int32
	AutoDepthStencilFormat     FORMAT
	Flags                      uint32
	FullScreen_RefreshRateInHz uint32
	PresentationInterval       uint32
}

func (obj *Direct3D) CreateDevice(
	adapter uint,
	deviceType DEVTYPE,
	focusWindow HWND,
	behaviorFlags uint32,
	params PRESENT_PARAMETERS,
) (*Device, PRESENT_PARAMETERS) {
	var device *Device
	syscall.Syscall9(
		obj.vtbl.CreateDevice,
		7,
		uintptr(unsafe.Pointer(obj)),
		uintptr(adapter),
		uintptr(deviceType),
		uintptr(focusWindow),
		uintptr(behaviorFlags),
		uintptr(unsafe.Pointer(&params)),
		uintptr(unsafe.Pointer(&device)),
		0,
		0,
	)
	return device, params
}

type deviceVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	TestCooperativeLevel        uintptr
	GetAvailableTextureMem      uintptr
	EvictManagedResources       uintptr
	GetDirect3D                 uintptr
	GetDeviceCaps               uintptr
	GetDisplayMode              uintptr
	GetCreationParameters       uintptr
	SetCursorProperties         uintptr
	SetCursorPosition           uintptr
	ShowCursor                  uintptr
	CreateAdditionalSwapChain   uintptr
	GetSwapChain                uintptr
	GetNumberOfSwapChains       uintptr
	Reset                       uintptr
	Present                     uintptr
	GetBackBuffer               uintptr
	GetRasterStatus             uintptr
	SetDialogBoxMode            uintptr
	SetGammaRamp                uintptr
	GetGammaRamp                uintptr
	CreateTexture               uintptr
	CreateVolumeTexture         uintptr
	CreateCubeTexture           uintptr
	CreateVertexBuffer          uintptr
	CreateIndexBuffer           uintptr
	CreateRenderTarget          uintptr
	CreateDepthStencilSurface   uintptr
	UpdateSurface               uintptr
	UpdateTexture               uintptr
	GetRenderTargetData         uintptr
	GetFrontBufferData          uintptr
	StretchRect                 uintptr
	ColorFill                   uintptr
	CreateOffscreenPlainSurface uintptr
	SetRenderTarget             uintptr
	GetRenderTarget             uintptr
	SetDepthStencilSurface      uintptr
	GetDepthStencilSurface      uintptr
	BeginScene                  uintptr
	EndScene                    uintptr
	Clear                       uintptr
	SetTransform                uintptr
	GetTransform                uintptr
	MultiplyTransform           uintptr
	SetViewport                 uintptr
	GetViewport                 uintptr
	SetMaterial                 uintptr
	GetMaterial                 uintptr
	SetLight                    uintptr
	GetLight                    uintptr
	LightEnable                 uintptr
	GetLightEnable              uintptr
	SetClipPlane                uintptr
	GetClipPlane                uintptr
	SetRenderState              uintptr
	GetRenderState              uintptr
	CreateStateBlock            uintptr
	BeginStateBlock             uintptr
	EndStateBlock               uintptr
	SetClipStatus               uintptr
	GetClipStatus               uintptr
	GetTexture                  uintptr
	SetTexture                  uintptr
	GetTextureStageState        uintptr
	SetTextureStageState        uintptr
	GetSamplerState             uintptr
}

type Device struct {
	vtbl *deviceVtbl
}

func (obj *Device) CreateOffscreenPlainSurface(
	width uint,
	height uint,
	format FORMAT,
	pool POOL,
	sharedHandle uintptr,
) *Surface {
	var surface *Surface
	syscall.Syscall9(
		obj.vtbl.CreateOffscreenPlainSurface,
		7,
		uintptr(unsafe.Pointer(obj)),
		uintptr(width),
		uintptr(height),
		uintptr(format),
		uintptr(pool),
		uintptr(unsafe.Pointer(&surface)),
		sharedHandle,
		0,
		0,
	)
	return surface
}

type Surface struct {
	vtbl *surfaceVtbl
}

type surfaceVtbl struct {
	QueryInterface uintptr
	AddRef         uintptr
	Release        uintptr

	GetDevice       uintptr
	SetPrivateData  uintptr
	GetPrivateData  uintptr
	FreePrivateData uintptr
	SetPriority     uintptr
	GetPriority     uintptr
	PreLoad         uintptr
	GetType         uintptr
	GetContainer    uintptr
	GetDesc         uintptr
	LockRect        uintptr
	UnlockRect      uintptr
	GetDC           uintptr
	ReleaseDC       uintptr
}

func (obj *Device) Release() uint32 {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return uint32(ret)
}

func (obj *Device) GetFrontBufferData(swapChain uint, destSurface *Surface) {
	syscall.Syscall(
		obj.vtbl.GetFrontBufferData,
		3,
		uintptr(unsafe.Pointer(obj)),
		uintptr(swapChain),
		uintptr(unsafe.Pointer(destSurface)),
	)
}

func (obj *Surface) LockRect(
	rect *RECT,
	flags uint32,
) (lockedRect LOCKED_RECT) {
	syscall.Syscall6(
		obj.vtbl.LockRect,
		4,
		uintptr(unsafe.Pointer(obj)),
		uintptr(unsafe.Pointer(&lockedRect)),
		uintptr(unsafe.Pointer(rect)),
		uintptr(flags),
		0,
		0,
	)
	return
}

type LOCKED_RECT struct {
	Pitch int32
	PBits uintptr
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func (obj *Surface) UnlockRect() {
	syscall.Syscall(
		obj.vtbl.UnlockRect,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
}

func (obj *Surface) Release() uint32 {
	ret, _, _ := syscall.Syscall(
		obj.vtbl.Release,
		1,
		uintptr(unsafe.Pointer(obj)),
		0,
		0,
	)
	return uint32(ret)
}

func DisplayCaptureStart(width, height int) {
	if started {
		return
	}

	version := 32
	dll := syscall.NewLazyDLL(DX_DLL)
	direct3DCreate9 := dll.NewProc("Direct3DCreate9")

	obj, _, _ := direct3DCreate9.Call(uintptr(version))
	d3d := (*Direct3D)(unsafe.Pointer(obj))
	defer d3d.Release()

	mode, _ := d3d.GetAdapterDisplayMode(0)
	device, _ := d3d.CreateDevice(
		0,
		DEVTYPE_HAL,
		0,
		CREATE_SOFTWARE_VERTEXPROCESSING,
		PRESENT_PARAMETERS{
			Windowed:         1,
			BackBufferCount:  1,
			BackBufferWidth:  mode.Width,
			BackBufferHeight: mode.Height,
			SwapEffect:       SWAPEFFECT_DISCARD,
		},
	)

	surface := device.CreateOffscreenPlainSurface(
		uint(mode.Width),
		uint(mode.Height),
		FMT_A8R8G8B8,
		POOL_SYSTEMMEM,
		0,
	)

	//defer surface.Release()
	img := image.NewRGBA(image.Rect(0, 0, int(mode.Width), int(mode.Height)))
	started = true
	go func() {
		for started {
			device.GetFrontBufferData(0, surface)
			r := surface.LockRect(nil, 0)
			surface.UnlockRect()

			if r.Pitch != int32(mode.Width*4) {
				continue
			}
			for i := range img.Pix {
				img.Pix[i] = *((*byte)(unsafe.Pointer(r.PBits + uintptr(i))))
			}
			for i := 0; i < len(img.Pix); i += 4 {
				img.Pix[i+0], img.Pix[i+2] = img.Pix[i+2], img.Pix[i+0]
			}

			callbackFrameWriter(img.Pix, 0)
		}
	}()
}

var started = false

func DisplayCaptureStop() {
	started = false
}

func GetDisplayRect() *DisplayRect {
	version := 32
	dll := syscall.NewLazyDLL(DX_DLL)
	direct3DCreate9 := dll.NewProc("Direct3DCreate9")
	obj, _, _ := direct3DCreate9.Call(uintptr(version))
	if obj == 0 {
		return nil
	}
	d3d := (*Direct3D)(unsafe.Pointer(obj))
	defer d3d.Release()
	mode, err := d3d.GetAdapterDisplayMode(0)
	if err != nil {
		return nil
	}

	return &DisplayRect{Width: int(mode.Width), Height: int(mode.Height)}
}
