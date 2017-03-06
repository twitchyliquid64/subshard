from ds_store import DSStore
# based on dmgbuild by al45tair, but without dependency hell

bounds = ((100, 100), (640, 280))

bwsp = {
    b'ShowStatusBar': False,
    b'WindowBounds': b'{{%s, %s}, {%s, %s}}' % (bounds[0][0],
                                                bounds[0][1],
                                                bounds[1][0],
                                                bounds[1][1]),
    b'ContainerShowSidebar': False,
    b'PreviewPaneVisibility': False,
    b'SidebarWidth': 180,
    b'ShowTabView': False,
    b'ShowToolbar': False,
    b'ShowPathbar': False,
    b'ShowSidebar': False
    }

icvp = {
    b'viewOptionsVersion': 1,
    b'backgroundType': 1,
    b'backgroundColorRed': 1.0,
    b'backgroundColorGreen': 0.9,
    b'backgroundColorBlue': 0.9,
    b'gridOffsetX': float(0),
    b'gridOffsetY': float(0),
    b'gridSpacing': float(100.0),
    b'arrangeBy': 'none',
    b'showIconPreview': False,
    b'showItemInfo': False,
    b'labelOnBottom': True,
    b'textSize': float(16.0),
    b'iconSize': float(128.0),
    b'scrollPositionX': float(0),
    b'scrollPositionY': float(0)
    }

icvl = (b'type', b'icnv')

def Make(output_path, icon_locations={}, include_icon_view_settings=True):
    background_file_present = False

    with DSStore.open(output_path, 'w+') as d:
            d['.']['vSrn'] = ('long', 1)
            d['.']['bwsp'] = bwsp
            if include_icon_view_settings: #optional
                d['.']['icvp'] = icvp
                if background_file_present:
                    d['.']['pBBk'] = background_bmk
            d['.']['icvl'] = icvl
            for k in icon_locations:
                d[k]['Iloc'] = icon_locations[k]
